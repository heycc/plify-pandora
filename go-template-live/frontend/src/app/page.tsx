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

  // Real WASM extraction function with fallback
  const extractVariables = useCallback(async (template: string): Promise<VariableInfo[]> => {
    try {
      const result = await wasmUtils.extractVariables(template);
      setError(null); // Clear error on success
      return Array.isArray(result) ? result : [];
    } catch (error) {
      const errorMessage = `Failed to extract variables: ${error instanceof Error ? error.message : 'Unknown error'}`;
      setError(errorMessage);
      console.error('Failed to extract variables with WASM:', error);
      return [];
    }
  }, [wasmLoaded]);

  // Real WASM rendering function
  const renderTemplate = useCallback(async (template: string, variables: Record<string, string>): Promise<string> => {
    try {
      const result = await wasmUtils.renderTemplate(template, variables);
      setError(null); // Clear error on success
      return typeof result === 'string' ? result : '';
    } catch (error) {
      const errorMessage = `Failed to render template: ${error instanceof Error ? error.message : 'Unknown error'}`;
      setError(errorMessage);
      console.error('Failed to render template with WASM:', error);
      return '';
    }
  }, [wasmLoaded]);

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
      const rendered = await renderTemplate(templateContent, variableValues);
      setRenderedContent(rendered);
    };

    render();
  }, [templateContent, variableValues, renderTemplate]);

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
        // Do not enable mock mode; keep WASM disabled on failure
        setWasmLoaded(false);
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
      <div className="mx-auto">
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
              <Button variant="outline" size="sm" onClick={() => loadRenderExample('profile')}>Profile</Button>
              <Button variant="outline" size="sm" onClick={() => loadRenderExample('config')}>Config</Button>
              <Button variant="outline" size="sm" onClick={() => loadRenderExample('email')}>Email</Button>
              <Button variant="outline" size="sm" onClick={clearTemplate}>🗑️ Clear All</Button>
            </div>
          </CardContent>
        </Card>


        {/* Main Editor Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-2 h-[calc(100vh-300px)] min-h-[600px]">
          {/* Left: Template Editor */}
          <div className="lg:col-span-3">
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