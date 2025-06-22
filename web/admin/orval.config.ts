import { defineConfig } from "orval";

export default defineConfig({
	adminApi: {
		input: {
			target: "../../typespec/tsp-output/@typespec/openapi3/openapi.yaml",
		},
		output: {
			mode: "single",
			target: "./src/generated-client/index.ts",
			schemas: "./src/generated-client/model",
			client: "fetch",
			override: {
				mutator: {
					path: "./src/api/custom-instance.ts",
					name: "customInstance",
				},
				useDates: false,
				useNamedParameters: true,
			},
		},
		hooks: {
			afterAllFilesWrite: "npm run format -- ./src/generated-client",
		},
	},
});
