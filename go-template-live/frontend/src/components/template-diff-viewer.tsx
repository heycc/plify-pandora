'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useEffect, useRef, useState } from 'react';
import type { WasmModel, WasmModelConfig } from '@/lib/wasm-utils';

interface TemplateDiffViewerProps {
  original: string;
  modified: string;
  onOriginalChange?: (value: string) => void;
  readOnly?: boolean;
  error?: string | null;
  wasmLoaded?: boolean;
  onShare?: () => void;
  shareStatus?: 'idle' | 'copying' | 'success' | 'error';
  selectedModel?: WasmModel;
  onModelChange?: (model: WasmModel) => void;
  availableModels?: Record<WasmModel, WasmModelConfig>;
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
  wasmLoaded = false,
  onShare,
  shareStatus = 'idle',
  selectedModel = 'official',
  onModelChange,
  availableModels
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
      theme: 'vs-oceanic',
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
    <Card className="h-full flex flex-col gap-2 py-4 bg-gray-100 backdrop-blur-sm">
      <CardHeader className="">
        <CardTitle className="text-lg flex items-center justify-between gap-2">
          <div className="flex items-center gap-3">
            <span className="text-2xl">📝</span>
            {/* Model Selector */}
            {availableModels && onModelChange && (
              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-2">
                  <select
                    value={selectedModel}
                    onChange={(e) => onModelChange(e.target.value as WasmModel)}
                    className="px-3 py-1.5 text-sm font-medium bg-white border-2 border-gray-300 rounded-lg shadow-sm hover:border-blue-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    disabled={!wasmLoaded}
                  >
                    {Object.entries(availableModels).map(([key, config]) => (
                      <option key={key} value={key}>
                        {config.name} - {config.description}
                      </option>
                    ))}
                  </select>
                  {!wasmLoaded && (
                    <span className="text-xs text-blue-600 animate-pulse">Loading...</span>
                  )}
                  {wasmLoaded && selectedModel && (
                    <span className="text-xs text-green-600 font-semibold">✓ Active</span>
                  )}
                </div>
                {selectedModel && availableModels[selectedModel] && (
                  <div className="text-xs text-gray-600 ml-1">
                    Functions: {availableModels[selectedModel].functions}
                  </div>
                )}
              </div>
            )}
          </div>
          <div className="flex items-center gap-2 ml-auto">
            {/* Error Display */}
            {error && (
              <div className="flex items-center gap-2 px-2 py-1 rounded-lg transition-all duration-300 bg-red-100 border-2 text-red-800 border-red-300 shadow-lg shadow-red-200/50"
                style={{
                  animation: 'pulse 2s cubic-bezier(0.4, 0, 0.8, 1) infinite',
                  boxShadow: '0 0 20px rgba(3, 1, 1, 0.3)'
                }}>
                <div className="w-5 h-5 flex items-center justify-center text-lg">
                  ⚠️
                </div>
                <div className="text-sm text-red-700">
                  {error}
                </div>
              </div>
            )}
            {/* WASM Loading Status */}
            {!wasmLoaded && !error && (
            <div className="flex items-center gap-2 py-1 px-2 rounded-lg transition-all duration-300 bg-blue-50 border border-blue-200 text-blue-700">
              <div className="w-5 h-5 flex items-center justify-center text-lg">
                ⏳
              </div>
              <div className="text-sm text-blue-600">
                Loading WebAssembly module...
              </div>
              </div>
            )}
            {/* Share Button */}
            {onShare && (
              <Button
                className="hover:cursor-pointer bg-sky-700 hover:bg-sky-800 text-white transition-all duration-200"
                variant="default"
                size="sm"
                onClick={onShare}
                disabled={shareStatus === 'copying' || !original.trim()}
              >
                {shareStatus === 'copying' ? '📋 Copying...' :
                 shareStatus === 'success' ? '✅ Copied!' :
                 '🔗 Share'}
              </Button>
            )}
          </div>
        </CardTitle>
      </CardHeader>

      <CardContent className="flex-1 px-6">
        <div className="h-full flex flex-col">
          {/* Monaco DiffEditor Container */}
          <div
            ref={containerRef}
            className="flex-1 min-h-[140px] py-0 border-8 border-gray-200/50 rounded-lg"
            style={{ height: '100%' }}
          />
          <div className="flex justify-between items-center px-4 py-2 mt-1 bg-muted text-sm text-muted-foreground">
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
        </div>
      </CardContent>
    </Card>
  );
}
