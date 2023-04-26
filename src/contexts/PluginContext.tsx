// create a context

import { FC, createContext, useContext, useEffect, useState } from "react";
import { PluginManager, IPlugin } from "@dottless/plugin";

export interface PluginContextI {
  plugins: IPlugin<any, any, any>[];
  loading: boolean;
}
export const PluginContext = createContext<PluginContextI>({
  plugins: [],
  loading: true,
});
export const usePluginContext = () => useContext(PluginContext);

export const PluginProvider: FC<any> = (props: any) => {
  const [loading, setLoading] = useState(true);
  const [plugins, setPlugins] = useState<any[]>([]);
  const pluginManager = PluginManager.getInstance();

  useEffect(() => {
    pluginManager
      .loadPlugins((window as any).__plugin_list_modules)
      .then((plugins) => {
        setPlugins(plugins);
        setLoading(false);
      })
      .catch((e) => {
        setLoading(false);
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (loading) return <div>Loading Plugins...</div>;

  return (
    <PluginContext.Provider value={{ plugins, loading }}>
      {props.children}
    </PluginContext.Provider>
  );
};
