'use client';

import { useState, useEffect, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { TemplateDiffViewer } from '@/components/template-diff-viewer';
import { VariablePanel } from '@/components/variable-panel';
import { examples, defaultExamples, renderExamples } from '@/lib/template-examples';
import { wasmUtils } from '@/lib/wasm-utils';

interface VariableInfo {
  name: string;
  defaultValue?: string;
}

export default function Home() {
  const [wasmLoaded, setWasmLoaded] = useState(false);

  // Template states
  const [templateContent, setTemplateContent] = useState(examples.simple);
  const [renderedContent, setRenderedContent] = useState('');

  // Variable states
  const [extractedVariables, setExtractedVariables] = useState<VariableInfo[]>([]);
  const [variableValues, setVariableValues] = useState<Record<string, string>>({});

  // Error state
  const [error, setError] = useState<string | null>(null);

  // Mock extraction functions (fallback if WASM fails)
  const extractVariablesMock = useCallback((template: string): VariableInfo[] => {
    const variables: VariableInfo[] = [];

    // Extract {{.variable}} patterns
    const fieldMatches = template.matchAll(/\{\{\s*\.([\w.]+)\s*\}\}/g);
    for (const match of fieldMatches) {
      variables.push({ name: match[1] });
    }

    // Extract getv("key", "default") patterns
    const getvMatches = template.matchAll(/getv\s*\(\s*"([^"]+)"\s*(?:,\s*"([^"]*)")?\s*\)/g);
    for (const match of getvMatches) {
      variables.push({
        name: match[1],
        defaultValue: match[2] || undefined
      });
    }

    // Remove duplicates
    return variables.filter((v, i, arr) =>
      arr.findIndex(item => item.name === v.name) === i
    );
  }, []);

  // Mock rendering function (fallback if WASM fails)
  const renderTemplateMock = useCallback((template: string, variables: Record<string, string>): string => {
    let result = template;

    // Replace {{.variable}} patterns
    result = result.replace(/\{\{\s*\.([\w.]+)\s*\}\}/g, (match, variableName) => {
      return variables[variableName] || '';
    });

    // Replace getv("key", "default") patterns
    result = result.replace(/getv\s*\(\s*"([^"]+)"\s*(?:,\s*"([^"]*)")?\s*\)/g, (match, key, defaultValue) => {
      return variables[key] || defaultValue || '';
    });

    // Simple condition handling for demo
    result = result.replace(/\{\{\s*if\s+\.([\w.]+)\s*\}\}([\s\S]*?)\{\{\s*end\s*\}\}/g, (match, condition, content) => {
      return variables[condition] ? content : '';
    });

    return result;
  }, []);

  // Real WASM extraction function with fallback
  const extractVariables = useCallback(async (template: string): Promise<VariableInfo[]> => {
    if (!wasmLoaded) {
      return extractVariablesMock(template);
    }

    try {
      const result = await wasmUtils.extractVariables(template);
      setError(null); // Clear error on success
      return result.length > 0 ? result : extractVariablesMock(template);
    } catch (error) {
      const errorMessage = `Failed to extract variables: ${error instanceof Error ? error.message : 'Unknown error'}`;
      setError(errorMessage);
      console.error('Failed to extract variables with WASM, using mock:', error);
      return extractVariablesMock(template);
    }
  }, [wasmLoaded, extractVariablesMock]);

  // Real WASM rendering function with fallback
  const renderTemplate = useCallback(async (template: string, variables: Record<string, string>): Promise<string> => {
    if (!wasmLoaded) {
      return renderTemplateMock(template, variables);
    }

    try {
      const result = await wasmUtils.renderTemplate(template, variables);
      setError(null); // Clear error on success
      return result || renderTemplateMock(template, variables);
    } catch (error) {
      const errorMessage = `Failed to render template: ${error instanceof Error ? error.message : 'Unknown error'}`;
      setError(errorMessage);
      console.error('Failed to render template with WASM, using mock:', error);
      return renderTemplateMock(template, variables);
    }
  }, [wasmLoaded, renderTemplateMock]);

  // Extract variables when template changes
  useEffect(() => {
    const extractVars = async () => {
      const variables = await extractVariables(templateContent);
      setExtractedVariables(variables);

      // Initialize variable values with defaults
      const initialValues: Record<string, string> = {};
      variables.forEach(variable => {
        initialValues[variable.name] = variable.defaultValue || '';
      });
      setVariableValues(initialValues);
    };

    extractVars();
  }, [templateContent, extractVariables]);

  // Re-render when template or variables change
  useEffect(() => {
    const render = async () => {
      if (wasmLoaded) {
        const rendered = await renderTemplate(templateContent, variableValues);
        setRenderedContent(rendered);
      }
    };

    render();
  }, [templateContent, variableValues, wasmLoaded, renderTemplate]);

  // Initialize WASM
  useEffect(() => {
    const initializeWASM = async () => {
      try {
        await wasmUtils.initialize();
        setWasmLoaded(true);
        setError(null); // Clear error on success
      } catch (error) {
        const errorMessage = `Failed to initialize WASM: ${error instanceof Error ? error.message : 'Unknown error'}`;
        setError(errorMessage);
        console.error('Failed to initialize WASM:', error);
        // Fallback to mock mode if WASM fails
        setWasmLoaded(true);
      }
    };

    initializeWASM();
  }, []);

  const handleVariableValuesChange = useCallback((values: Record<string, string>) => {
    setVariableValues(values);
  }, []);

  const loadExample = (type: keyof typeof examples) => {
    setTemplateContent(examples[type]);
  };

  const loadDefaultExample = (type: keyof typeof defaultExamples) => {
    setTemplateContent(defaultExamples[type]);
  };

  const loadRenderExample = (type: keyof typeof renderExamples) => {
    setTemplateContent(renderExamples[type]);
  };

  const clearTemplate = () => {
    setTemplateContent('');
    setExtractedVariables([]);
    setVariableValues({});
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-600 to-purple-700 p-4">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <Card className="mb-2 border-0 shadow-xl gap-2">
          <CardHeader className="text-center bg-gradient-to-r from-blue-600 to-purple-700 text-white">
            <CardTitle className="text-3xl font-bold mb-2">🚀 Go Template Parser Demo</CardTitle>
            <CardDescription className="text-blue-100 text-lg">
              Interactive template editor with real-time variable extraction and rendering
            </CardDescription>
          </CardHeader>
        </Card>

        {/* Example Templates */}
        <Card className="mb-2">
          <CardHeader>
            <CardTitle className="text-lg">Quick Start Examples</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex gap-2 flex-wrap">
              <Button variant="outline" size="sm" onClick={() => loadExample('simple')}>Simple</Button>
              <Button variant="outline" size="sm" onClick={() => loadExample('complex')}>Complex</Button>
              <Button variant="outline" size="sm" onClick={() => loadExample('functions')}>Functions</Button>
              <Button variant="outline" size="sm" onClick={() => loadExample('control')}>Control Flow</Button>
              <Button variant="outline" size="sm" onClick={() => loadDefaultExample('basic')}>Basic Defaults</Button>
              <Button variant="outline" size="sm" onClick={() => loadDefaultExample('advanced')}>Advanced Defaults</Button>
              <Button variant="outline" size="sm" onClick={() => loadRenderExample('profile')}>Profile</Button>
              <Button variant="outline" size="sm" onClick={() => loadRenderExample('config')}>Config</Button>
              <Button variant="outline" size="sm" onClick={() => loadRenderExample('email')}>Email</Button>
              <Button variant="outline" size="sm" onClick={clearTemplate}>🗑️ Clear All</Button>
            </div>
          </CardContent>
        </Card>


        {/* Main Editor Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-2 h-[calc(100vh-300px)] min-h-[600px]">
          {/* Left: Template Editor */}
          <div className="lg:col-span-2">
            <TemplateDiffViewer
              original={templateContent}
              modified={renderedContent}
              onOriginalChange={setTemplateContent}
              readOnly={!wasmLoaded}
              error={error}
              wasmLoaded={wasmLoaded}
            />
          </div>

          {/* Right: Variables Panel */}
          <div className="lg:col-span-1">
            <VariablePanel
              variables={extractedVariables}
              onVariablesChange={handleVariableValuesChange}
            />
          </div>
        </div>

        {/* Info Footer */}
        <Card className="mt-6">
          <CardContent className="pt-6">
            <div className="text-center text-sm text-muted-foreground">
              <p>
                <strong>Supported Syntax:</strong> Go template syntax with variables, functions, and control structures
              </p>
              <p className="mt-1">
                Edit the template on the left, see extracted variables on the right, and watch the real-time rendering in the diff viewer.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}