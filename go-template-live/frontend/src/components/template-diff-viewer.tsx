'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useEffect, useRef, useState } from 'react';

interface TemplateDiffViewerProps {
  original: string;
  modified: string;
  onOriginalChange?: (value: string) => void;
  readOnly?: boolean;
  error?: string | null;
  wasmLoaded?: boolean;
}

declare global {
  interface Window {
    require: any;
    MonacoEnvironment?: any;
  }
}

export function TemplateDiffViewer({
  original,
  modified,
  onOriginalChange,
  readOnly = false,
  error = null,
  wasmLoaded = false
}: TemplateDiffViewerProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const diffEditorRef = useRef<any>(null);
  const monacoRef = useRef<any>(null);
  const originalModelRef = useRef<any>(null);
  const modifiedModelRef = useRef<any>(null);
  const [isMonacoLoaded, setIsMonacoLoaded] = useState(false);

  // Load Monaco Editor from CDN
  useEffect(() => {
    // Check if Monaco is already loaded
    if (window.require && window.require.defined && window.require.defined('vs/editor/editor.main')) {
      window.require(['vs/editor/editor.main'], (monaco: any) => {
        monacoRef.current = monaco;
        setIsMonacoLoaded(true);
      });
      return;
    }

    // Load Monaco Editor loader
    const loaderScript = document.createElement('script');
    loaderScript.src = 'https://unpkg.com/monaco-editor@0.45.0/min/vs/loader.js';
    loaderScript.async = true;

    loaderScript.onload = () => {
      // Configure Monaco Environment
      window.MonacoEnvironment = {
        getWorkerUrl: function (_: any, label: string) {
          if (label === 'json') {
            return 'https://unpkg.com/monaco-editor@0.45.0/min/vs/language/json/json.worker.js';
          }
          if (label === 'css' || label === 'scss' || label === 'less') {
            return 'https://unpkg.com/monaco-editor@0.45.0/min/vs/language/css/css.worker.js';
          }
          if (label === 'html' || label === 'handlebars' || label === 'razor') {
            return 'https://unpkg.com/monaco-editor@0.45.0/min/vs/language/html/html.worker.js';
          }
          if (label === 'typescript' || label === 'javascript') {
            return 'https://unpkg.com/monaco-editor@0.45.0/min/vs/language/typescript/ts.worker.js';
          }
          return 'https://unpkg.com/monaco-editor@0.45.0/min/vs/base/worker/workerMain.js';
        }
      };

      // Configure loader
      window.require.config({
        paths: {
          vs: 'https://unpkg.com/monaco-editor@0.45.0/min/vs'
        }
      });

      // Load Monaco Editor
      window.require(['vs/editor/editor.main'], (monaco: any) => {
        monacoRef.current = monaco;
        setIsMonacoLoaded(true);
      });
    };

    document.head.appendChild(loaderScript);

    return () => {
      // Cleanup is handled by Monaco itself
    };
  }, []);

  // Initialize DiffEditor
  useEffect(() => {
    if (!isMonacoLoaded || !monacoRef.current || !containerRef.current || diffEditorRef.current) {
      return;
    }

    const monaco = monacoRef.current;

    // Create models
    originalModelRef.current = monaco.editor.createModel(original, 'html');
    modifiedModelRef.current = monaco.editor.createModel(modified, 'html');

    // Create diff editor
    const diffEditor = monaco.editor.createDiffEditor(containerRef.current, {
      enableSplitViewResizing: true,
      renderSideBySide: true,
      originalEditable: !readOnly,
      readOnly: true,
      automaticLayout: true,
      scrollBeyondLastLine: false,
      minimap: { enabled: false },
      fontSize: 13,
      lineNumbers: 'on',
      lineNumbersMinChars: 3,
      wordWrap: 'on',
      theme: 'vs-dark',
      diffWordWrap: 'on',
      ignoreTrimWhitespace: false,
      renderIndicators: true,
      diffAlgorithm: 'advanced',
      renderSideBySideInlineBreakpoint: 0, // Disable automatic inline mode switching
    });

    // Set models
    diffEditor.setModel({
      original: originalModelRef.current,
      modified: modifiedModelRef.current,
    });

    diffEditorRef.current = diffEditor;

    // Listen for changes in the original editor (left side)
    if (onOriginalChange && !readOnly) {
      const originalEditor = diffEditor.getOriginalEditor();
      originalEditor.onDidChangeModelContent(() => {
        const newValue = originalEditor.getValue();
        if (onOriginalChange) {
          onOriginalChange(newValue);
        }
      });
    }

    // Cleanup
    return () => {
      if (diffEditorRef.current) {
        diffEditorRef.current.dispose();
        diffEditorRef.current = null;
      }
      if (originalModelRef.current) {
        originalModelRef.current.dispose();
        originalModelRef.current = null;
      }
      if (modifiedModelRef.current) {
        modifiedModelRef.current.dispose();
        modifiedModelRef.current = null;
      }
    };
  }, [isMonacoLoaded, readOnly, onOriginalChange]);

  // Update content when props change
  useEffect(() => {
    if (!originalModelRef.current || !modifiedModelRef.current) return;

    const currentOriginal = originalModelRef.current.getValue();
    const currentModified = modifiedModelRef.current.getValue();

    // Only update if different to avoid cursor jumps
    if (currentOriginal !== original) {
      originalModelRef.current.setValue(original);
    }

    if (currentModified !== modified) {
      modifiedModelRef.current.setValue(modified);
    }
  }, [original, modified]);

  return (
    <Card className="h-full flex flex-col gap-2 py-4">
      <CardHeader className="">
        <CardTitle className="text-lg flex items-center gap-2">
          <span className="text-2xl">📝</span>
          Template Diff Viewer
        </CardTitle>
      </CardHeader>

      <CardContent className="flex-1 px-6">
        <div className="h-full flex flex-col">
          {/* Monaco DiffEditor Container */}
          <div
            ref={containerRef}
            className="flex-1 min-h-[140px] bg-gray-900 py-2 border-2 border-gray-800 rounded"
            style={{ height: '100%' }}
          />
          <div className="flex justify-between items-center px-4 py-2 bg-muted text-sm text-muted-foreground border-t">
            <div className="text-xs">
              {readOnly ? 'Read-only Diff View' : '✏️ Edit template on left to see live rendering on right'}
            </div>
            <div className="text-xs flex items-center gap-4 ">
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded" style={{ backgroundColor: 'rgba(255, 100, 100, 0.3)' }}></div>
                <span>Removed/Changed</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded" style={{ backgroundColor: 'rgba(100, 255, 100, 0.3)' }}></div>
                <span>Added/New</span>
              </div>
            </div>
          </div>
          {/* Error Display */}
          {error && (
            <div className="mt-3">
              <div className="flex items-center gap-3 p-2 rounded-lg transition-all duration-300 bg-red-100 border-2 text-red-800 border-red-300 shadow-lg shadow-red-200/50"
                style={{
                  animation: 'pulse 2s cubic-bezier(0.4, 0, 0.8, 1) infinite',
                  boxShadow: '0 0 20px rgba(3, 1, 1, 0.3)'
                }}>
                <div className="w-6 h-6 flex items-center justify-center text-xl">
                  ⚠️
                </div>
                <div className="flex-1">
                  <div className="text-sm mt-1 text-red-700">
                    {error}
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* WASM Loading Status */}
          {!wasmLoaded && !error && (
            <div className="mt-3">
              <div className="flex items-center gap-3 p-2 rounded-lg transition-all duration-300 bg-blue-50 border border-blue-200 text-blue-700">
                <div className="w-6 h-6 flex items-center justify-center text-xl">
                  ⏳
                </div>
                <div className="flex-1">
                  <div className="text-sm mt-1 text-blue-600">
                    Loading WebAssembly module...
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>

      </CardContent>
    </Card>
  );
}
