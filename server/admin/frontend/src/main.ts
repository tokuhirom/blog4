import { mount } from "svelte";
import "./app.css";
import App from "./App.svelte";

const appElement = document.getElementById("app");
if (!appElement) {
	throw new Error("Element with id 'app' not found");
}
const app = mount(App, {
	target: appElement,
});

export default app;
