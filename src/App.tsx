import { FC } from "react";
import { PluginProvider } from "./contexts/PluginContext";
import { PluginsList } from "./pages/PluginsList";

const App: FC = () => {
  return (
    <PluginProvider>
      <div className="App">Main APP</div>
      <PluginsList />
    </PluginProvider>
  );
};

export default App;
