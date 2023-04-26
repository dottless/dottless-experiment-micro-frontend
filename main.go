package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Prepare the index file
	if err := prepareIndex(); err != nil {
		panic(err)
	}

	build := http.FileServer(http.Dir("build"))
	handler := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/plugins") {
			// Get the file path by replacing "/plugins" with "plugins" and removing any trailing slashes
			filePath := strings.TrimPrefix(r.URL.Path, "/plugins")
			filePath = strings.TrimPrefix(filePath, "/")
			filePath = filepath.Join("plugins", filePath)

			// Check if the file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				fmt.Println("File not found: " + filePath)
				http.NotFound(w, r)
				return
			}
			http.ServeFile(w, r, filePath)
		} else {
			build.ServeHTTP(w, r)
		}
	}

	// Register the request handler function
	http.HandleFunc("/", handler)

	// Start the server
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server started on port 3000")
}

type PluginsModule struct {
	Id   string `json:"id"`
	Path string `json:"modulePath"`
}

type IndexData struct {
	Plugins []PluginsModule `json:"Plugins"`
}

func prepareIndex() error {
	// Read the template file
	templateData, err := os.ReadFile("build/index-template.html")
	if err != nil {
		fmt.Errorf("Error reading template file: %s", err)
	}

	// Parse the template
	tmpl, err := template.New("index").Parse(string(templateData))
	if err != nil {
		return fmt.Errorf("Error parsing template: %s", err)
	}

	// Create the data for the page
	plugins, _ := preparePluginList()

	// Create a new file to write the output to
	outputFile, err := os.Create("build/index.html")
	if err != nil {
		return fmt.Errorf("Error creating output file: %s", err)
	}
	defer outputFile.Close()

	// Execute the template with the data and write the output to the file
	err = tmpl.Execute(outputFile, IndexData{
		Plugins: plugins,
	})

	if err != nil {
		return fmt.Errorf("Error executing template: %s", err)
	}

	return nil
}

func preparePluginList() ([]PluginsModule, error) {
	pluginsDir := "plugins"
	topLevelDirs := make(map[string]bool)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if path == pluginsDir {
			return nil
		}
		dir := filepath.Base(path)
		if _, ok := topLevelDirs[dir]; !ok {
			topLevelDirs[dir] = true
		}
		return filepath.SkipDir
	}
	if err := filepath.Walk(pluginsDir, walkFunc); err != nil {
		return nil, fmt.Errorf("Error walking directory: %v\n", err)
	}

	plugins := make([]PluginsModule, 0)
	// Print the top-level directories
	for dir := range topLevelDirs {
		plugin, err := getPluginMeta(dir)
		if err != nil {
			fmt.Println(err)
			continue
		}
		plugins = append(plugins, PluginsModule{
			Id:   plugin.Id,
			Path: plugin.Path,
		})
	}

	return plugins, nil
}

type Plugin struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Path        string `json:"-"`
}

func getPluginMeta(folderName string) (*Plugin, error) {
	// Read the contents of the file
	pluginPath := filepath.Join("plugins", folderName, "dist", "plugin.json")
	fileBytes, err := ioutil.ReadFile(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %s", err)
	}

	// Unmarshal the JSON data into a Plugin struct
	var plugin Plugin
	err = json.Unmarshal(fileBytes, &plugin)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s", err)
	}

	plugin.Path = filepath.Join("plugins", folderName)
	return &plugin, nil
}
