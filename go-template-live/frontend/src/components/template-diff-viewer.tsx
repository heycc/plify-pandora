'use client';

import { useState, useEffect, useRef } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

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
      readOnly: readOnly,
      automaticLayout: true,
      scrollBeyondLastLine: false,
      minimap: { enabled: false },
      fontSize: 14,
      lineNumbers: 'on',
      wordWrap: 'on',
      theme: 'vs-dark',
      diffWordWrap: 'on',
      ignoreTrimWhitespace: false,
      renderIndicators: true,
      diffAlgorithm: 'advanced',
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
    <Card className="h-full flex flex-col">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center gap-2">
          <span className="text-2xl">📝</span>
          Template Diff Viewer
        </CardTitle>
      </CardHeader>

      {/* Status & Error Display */}
      {(error || !wasmLoaded) && (
        <div className="px-6 pb-3">
          <div className={`flex items-center gap-3 p-3 rounded-lg ${error ? 'bg-red-50 border border-red-200 text-red-700' : 'bg-blue-50 border border-blue-200 text-blue-700'}`}>
            <div className={`w-6 h-6 rounded-full flex items-center justify-center text-sm ${
              error ? 'bg-red-600' : 'bg-blue-600'
            }`}>
              {error ? '⚠️' : '⏳'}
            </div>
            <div>
              <div className="font-medium">{error ? 'Error' : 'Status'}</div>
              <div className="text-sm">
                {error ? error : (
                  wasmLoaded ? '✅ WebAssembly module loaded successfully!' : 'Loading WebAssembly module...'
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      <CardContent className="flex-1 p-0">
        <div className="h-full flex flex-col">
          {/* Monaco DiffEditor Container */}
          <div 
            ref={containerRef} 
            className="flex-1 min-h-[400px]"
            style={{ height: '100%' }}
          />
          
          <div className="flex justify-between items-center px-4 py-2 bg-muted text-sm text-muted-foreground border-t">
            <div className="text-xs">
              {readOnly ? 'Read-only Diff View' : '✏️ Edit template on left to see live rendering on right'}
            </div>
            <div className="text-xs flex items-center gap-4">
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
