'use client';

import { useState, useEffect, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { TemplateDiffViewer } from '@/components/template-diff-viewer';
import { VariablePanel } from '@/components/variable-panel';
import { examples } from '@/lib/template-examples';
import { wasmUtils } from '@/lib/wasm-utils';
import { getTemplateTextFromUrl, copyShareableUrl, hasSharedTemplateInUrl, generateShareableUrl } from '@/lib/url-sharing';

interface VariableInfo {
  name: string;
  defaultValue?: string;
}

export default function Home() {
  const [wasmLoaded, setWasmLoaded] = useState(false);

  // Template states
  const [templateContent, setTemplateContent] = useState(examples.Basic);
  const [renderedContent, setRenderedContent] = useState('');

  // Variable states
  const [extractedVariables, setExtractedVariables] = useState<VariableInfo[]>([]);
  const [variableValues, setVariableValues] = useState<Record<string, string>>({});

  // Error state
  const [error, setError] = useState<string | null>(null);

  // Share state
  const [shareStatus, setShareStatus] = useState<'idle' | 'copying' | 'success' | 'error'>('idle');

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

  // Load template from URL if present
  useEffect(() => {
    if (hasSharedTemplateInUrl()) {
      const sharedTemplate = getTemplateTextFromUrl();
      if (sharedTemplate) {
        setTemplateContent(sharedTemplate);
      }
    }
  }, []);

  // Handle Cmd+S / Ctrl+S to update URL
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key === 's') {
        event.preventDefault(); // Prevent default browser save dialog
        
        if (!templateContent.trim()) {
          console.log('Cannot update URL with empty template');
          return;
        }

        try {
          const shareableUrl = generateShareableUrl(templateContent);
          if (shareableUrl) {
            // Update URL without reloading the page
            const url = new URL(shareableUrl);
            window.history.replaceState({}, '', url.toString());
            console.log('URL updated successfully');
          } else {
            console.warn('Template is too large to share via URL');
          }
        } catch (error) {
          console.error('Failed to update URL:', error);
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [templateContent]);

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

  // Helper function to get display name from key
  const getDisplayName = (key: string): string => {
    return key.replace(/_/g, ' ');
  };

  const clearTemplate = () => {
    setTemplateContent('');
    setExtractedVariables([]);
    setVariableValues({});
  };

  const handleShareTemplate = async () => {
    if (!templateContent.trim()) {
      setError('Cannot share empty template');
      return;
    }

    setShareStatus('copying');
    try {
      const success = await copyShareableUrl(templateContent);
      if (success) {
        setShareStatus('success');

        // Update current page URL to match the shared template
        const shareableUrl = generateShareableUrl(templateContent);
        if (shareableUrl) {
          // Update URL without reloading the page
          const url = new URL(shareableUrl);
          window.history.replaceState({}, '', url.toString());
        }

        // Reset success status after 2 seconds
        setTimeout(() => setShareStatus('idle'), 2000);
      } else {
        setShareStatus('error');
        setError('Template is too large to share via URL. The maximum allowed is about 2000 characters.');
      }
    } catch (error) {
      setShareStatus('error');
      setError('Failed to share template');
      console.error('Failed to share template:', error);
    }
  };

  return (
    <div className="h-screen mx-auto p-4 bg-gradient-to-b from-sky-200 to-amber-100 flex flex-col overflow-hidden">
      {/* Header */}
      <div className="mb-2 border-0 flex-shrink-0 py-2">
        <div className="text-center text-blue-700 py-2">
          <h1 className="text-3xl font-bold mb-2">Golang Template Live Preview</h1>
          <h2 className="text-gray-500 text-lg">
            Type in Golang template, extract variables, render and view the diff in real-time 🚀
          </h2>
        </div>
      </div>

      {/* Example Templates */}
      <section className="flex-shrink-0">
        <div className="flex flex-row gap-2 px-6 py-2 bg-white/50 backdrop-blur-sm rounded-lg items-center justify-start">
          <div className="py-0">
            <h3 className="text-medium font-semibold">Quick Start Examples</h3>
          </div>
          <div className="py-0">
            <div className="flex space-x-2 flex-wrap">
              {/* All Examples */}
              {Object.keys(examples).map(key => (
                <Button
                  key={key}
                  className="hover:cursor-pointer"
                  variant="outline"
                  size="sm"
                  onClick={() => loadExample(key as keyof typeof examples)}
                >
                  {getDisplayName(key)}
                </Button>
              ))}

              <Button className="hover:cursor-pointer" variant="outline" size="sm" onClick={clearTemplate}>🗑️ Clear All</Button>
            </div>
          </div>
        </div>
      </section>

      {/* Main Editor Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-4 gap-2 flex-1 min-h-0 mt-2">
        {/* Left: Template Editor */}
        <div className="lg:col-span-3 min-h-0">
          <TemplateDiffViewer
            original={templateContent}
            modified={renderedContent}
            onOriginalChange={setTemplateContent}
            readOnly={!wasmLoaded}
            error={error}
            wasmLoaded={wasmLoaded}
            onShare={handleShareTemplate}
            shareStatus={shareStatus}
          />
        </div>

        {/* Right: Variables Panel */}
        <div className="lg:col-span-1 min-h-0">
          <VariablePanel
            variables={extractedVariables}
            onVariablesChange={handleVariableValuesChange}
          />
        </div>
      </div>

      
    </div>
  );
}