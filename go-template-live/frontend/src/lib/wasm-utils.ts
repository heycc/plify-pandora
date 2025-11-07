// WASM utility functions for template parsing

export type WasmModel = 'official' | 'custom' | 'confd';

export interface WasmModelConfig {
  name: string;
  fileName: string;
  description: string;
  functions: string;
}

export const WASM_MODELS: Record<WasmModel, WasmModelConfig> = {
  official: {
    name: 'Official',
    fileName: 'official.wasm',
    description: 'Official Go template functions only',
    functions: 'Standard Go template functions (no custom functions)'
  },
  custom: {
    name: 'Custom',
    fileName: 'custom.wasm',
    description: 'Official + custom functions',
    functions: 'getv, exists, get, json, jsonArray'
  },
  confd: {
    name: 'Confd',
    fileName: 'confd.wasm',
    description: 'Official + Confd-style functions',
    functions: 'base, split, json, jsonArray, dir, map, join, datetime, toUpper, toLower, replace, contains, base64Encode, base64Decode, trimSuffix, parseBool, reverse, add, sub, div, mod, mul, seq, atoi'
  }
};

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
  private currentModel: WasmModel = 'official';
  private goInstance: any = null;

  async initialize(model: WasmModel = 'official'): Promise<void> {
    console.log(`[WASMUtils] Initialize called with model: ${model}, current: ${this.currentModel}, initialized: ${this.isInitialized}`);
    
    // If model changed, force re-initialization
    if (this.isInitialized && this.currentModel !== model) {
      console.log(`[WASMUtils] Model changed from ${this.currentModel} to ${model}, forcing re-initialization`);
      this.isInitialized = false;
      this.initializationPromise = null;
      this.goInstance = null;
    }

    if (this.isInitialized && this.currentModel === model) {
      console.log(`[WASMUtils] Already initialized with model ${model}`);
      return;
    }

    if (this.initializationPromise) {
      console.log(`[WASMUtils] Initialization in progress, waiting...`);
      return this.initializationPromise;
    }

    this.currentModel = model;
    this.initializationPromise = this.loadWASM(model);
    await this.initializationPromise;
    this.isInitialized = true;
    console.log(`[WASMUtils] Initialization complete for model: ${model}`);
  }

  getAvailableModels(): Record<WasmModel, WasmModelConfig> {
    return WASM_MODELS;
  }

  getCurrentModel(): WasmModel {
    return this.currentModel;
  }

  private async loadWASM(model: WasmModel): Promise<void> {
    try {
      const modelConfig = WASM_MODELS[model];
      console.log(`[WASMUtils] Starting WASM loading for model: ${modelConfig.name} (${modelConfig.fileName})...`);

      // Load the Go WASM runtime (only once)
      if (!(window as any).Go) {
        console.log('[WASMUtils] Loading wasm_exec.js...');
        await this.loadScript(`${BASE_PATH}/wasm_exec.js`);
        console.log('[WASMUtils] wasm_exec.js loaded');
      } else {
        console.log('[WASMUtils] wasm_exec.js already loaded');
      }

      // Create a new Go instance for this model
      const go = new (window as any).Go();
      this.goInstance = go;
      console.log('[WASMUtils] New Go instance created');

      // Fetch and instantiate the WASM module
      const wasmResponse = await fetch(`${BASE_PATH}/${modelConfig.fileName}`);
      console.log(`[WASMUtils] WASM fetch response for ${modelConfig.fileName}:`, wasmResponse.status);

      if (!wasmResponse.ok) {
        throw new Error(`Failed to fetch ${modelConfig.fileName}: ${wasmResponse.status}`);
      }

      const result = await WebAssembly.instantiateStreaming(
        wasmResponse,
        go.importObject
      );
      console.log('[WASMUtils] WASM instantiated successfully');

      // Run the Go program (this registers the global functions)
      go.run(result.instance);
      console.log('[WASMUtils] Go runtime started');

      // Wait a bit for functions to be registered
      await new Promise(resolve => setTimeout(resolve, 100));

      // Verify WASM functions are registered
      console.log('[WASMUtils] WASM functions registered:', {
        extractTemplateVariables: typeof window.extractTemplateVariables,
        extractTemplateVariablesSimple: typeof window.extractTemplateVariablesSimple,
        renderTemplateWithValues: typeof window.renderTemplateWithValues
      });

      if (typeof window.extractTemplateVariables !== 'function') {
        throw new Error('WASM functions not properly registered');
      }

      console.log(`[WASMUtils] ✓ Model ${modelConfig.name} loaded successfully`);
    } catch (error) {
      console.error('[WASMUtils] Failed to load WASM:', error);
      this.isInitialized = false;
      this.initializationPromise = null;
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
    // Don't call initialize() here - it's handled by the page component
    // Calling it without model parameter would reset to 'official'
    
    try {
      console.log(`[WASMUtils] Extracting variables with model: ${this.currentModel}`);
      console.log(`[WASMUtils] Function available:`, typeof window.extractTemplateVariables);
      
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
    // Don't call initialize() here - it's handled by the page component
    
    try {
      if (!window.extractTemplateVariablesSimple) {
        throw new Error('WASM function extractTemplateVariablesSimple not available');
      }
      
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
    // Don't call initialize() here - it's handled by the page component
    
    try {
      console.log(`[WASMUtils] Rendering template with model: ${this.currentModel}`);
      console.log(`[WASMUtils] Function available:`, typeof window.renderTemplateWithValues);
      
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