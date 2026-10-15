import { mount } from "svelte";
import "@kenn-io/kit-ui/theme.css";
import App from "./App.svelte";
import "./app.css";

mount(App, { target: document.getElementById("app")! });
