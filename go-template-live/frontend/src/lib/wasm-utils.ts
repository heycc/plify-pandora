// WASM utility functions for template parsing

interface VariableInfo {
  name: string;
  defaultValue?: string;
}

declare global {
  interface Window {
    extractTemplateVariables: (template: string, fileName?: string) => string | { error: string };
    extractTemplateVariablesSimple: (template: string, fileName?: string) => string | { error: string };
    renderTemplateWithValues: (template: string, variablesJSON: string) => string | { error: string };
  }
}

// Get the base path from the environment variable or default to empty string
const BASE_PATH = process.env.NEXT_PUBLIC_BASE_PATH || '';

class WASMUtils {
  private isInitialized = false;
  private initializationPromise: Promise<void> | null = null;

  async initialize(): Promise<void> {
    if (this.isInitialized) {
      return;
    }

    if (this.initializationPromise) {
      return this.initializationPromise;
    }

    this.initializationPromise = this.loadWASM();
    await this.initializationPromise;
    this.isInitialized = true;
  }

  private async loadWASM(): Promise<void> {
    try {
      console.log('Starting WASM loading...');

      // Load the Go WASM runtime
      await this.loadScript(`${BASE_PATH}/wasm_exec.js`);
      console.log('WASM runtime loaded');

      // Load and instantiate the WASM module
      const go = new (window as any).Go();
      console.log('Go instance created');

      const wasmResponse = await fetch(`${BASE_PATH}/main.wasm`);
      console.log('WASM fetch response:', wasmResponse.status);

      const result = await WebAssembly.instantiateStreaming(
        wasmResponse,
        go.importObject
      );
      console.log('WASM instantiated');

      go.run(result.instance);
      console.log('Go runtime started');

      // Verify WASM functions are registered
      console.log('WASM functions registered:', {
        extractTemplateVariables: typeof window.extractTemplateVariables,
        extractTemplateVariablesSimple: typeof window.extractTemplateVariablesSimple,
        renderTemplateWithValues: typeof window.renderTemplateWithValues
      });
    } catch (error) {
      console.warn('Failed to load WASM:', error);
      throw error;
    }
  }

  private async loadScript(src: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const script = document.createElement('script');
      script.src = src;
      script.onload = () => resolve();
      script.onerror = () => reject(new Error(`Failed to load script: ${src}`));
      document.head.appendChild(script);
    });
  }

  async extractVariables(template: string): Promise<VariableInfo[]> {
    await this.initialize();

    try {
      if (!window.extractTemplateVariables) {
        throw new Error('WASM function extractTemplateVariables not available');
      }

      const result = window.extractTemplateVariables(template, 'template.tmpl');

      // Handle error case
      if (result && typeof result === 'object' && 'error' in result) {
        throw new Error(result.error);
      }

      // Handle success case - result is a JSON string
      if (typeof result === 'string') {
        const parsed = JSON.parse(result);
        if (Array.isArray(parsed)) {
          // Deduplicate variables by name, keeping the first occurrence
          const seen = new Map<string, VariableInfo>();
          for (const variable of parsed) {
            if (!seen.has(variable.name)) {
              seen.set(variable.name, variable);
            }
          }
          return Array.from(seen.values());
        }
        return [];
      }

      return [];
    } catch (error) {
      console.warn('Failed to extract variables:', error);
      throw error;
    }
  }

  async extractVariablesSimple(template: string): Promise<string[]> {
    await this.initialize();

    try {
      const result = window.extractTemplateVariablesSimple(template, 'template.tmpl');

      // Handle error case
      if (result && typeof result === 'object' && 'error' in result) {
        throw new Error(result.error);
      }

      // Handle success case - result is a JSON string
      if (typeof result === 'string') {
        const parsed = JSON.parse(result);
        return Array.isArray(parsed) ? parsed : [];
      }

      return [];
    } catch (error) {
      console.warn('Failed to extract simple variables:', error);
      throw error;
    }
  }

  async renderTemplate(template: string, variables: Record<string, string>): Promise<string> {
    await this.initialize();

    try {
      if (!window.renderTemplateWithValues) {
        throw new Error('WASM function renderTemplateWithValues not available');
      }

      // Convert all values to strings to handle boolean and other types
      const stringVariables: Record<string, string> = {};
      for (const [key, value] of Object.entries(variables)) {
        stringVariables[key] = String(value);
      }

      const variablesJSON = JSON.stringify(stringVariables);
      const result = window.renderTemplateWithValues(template, variablesJSON);

      // Handle error case
      if (result && typeof result === 'object' && 'error' in result) {
        throw new Error(result.error);
      }

      // Handle success case - result is a string
      if (typeof result === 'string') {
        return result;
      }

      return '';
    } catch (error) {
      console.warn('Failed to render template:', error);
      // Provide more helpful error message for boolean values
      if (error instanceof Error && (error.message.includes('expected string; found true') || error.message.includes('expected string; found false'))) {
        throw new Error('Template contains boolean values (true/false) which are not supported. Use string values instead.');
      }
      throw new Error(error instanceof Error ? error.message : 'Failed to render template');
    }
  }
}

export const wasmUtils = new WASMUtils();