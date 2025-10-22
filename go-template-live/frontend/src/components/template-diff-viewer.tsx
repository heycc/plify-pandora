'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { DiffEditor } from '@monaco-editor/react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface TemplateDiffViewerProps {
  original: string;
  modified: string;
  onOriginalChange?: (value: string) => void;
  readOnly?: boolean;
  error?: string | null;
  wasmLoaded?: boolean;
}

export function TemplateDiffViewer({
  original,
  modified,
  onOriginalChange,
  readOnly = false,
  error = null,
  wasmLoaded = false
}: TemplateDiffViewerProps) {
  const diffEditorRef = useRef<any>(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [localOriginal, setLocalOriginal] = useState(original);

  // Debounced change handler
  const handleEditorChange = useCallback((value: string | undefined) => {
    if (onOriginalChange && value !== undefined) {
      // Clear any existing timeout
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      // Set new timeout for debounced update
      timeoutRef.current = setTimeout(() => {
        console.log('Editor change detected:', value);
        setLocalOriginal(value);
        onOriginalChange(value);
      }, 300); // 300ms debounce
    }
  }, [onOriginalChange]);

  // Sync local state with props
  useEffect(() => {
    if (original !== localOriginal) {
      setLocalOriginal(original);
    }
  }, [original]);

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  // Update editor content when props change
  useEffect(() => {
    if (diffEditorRef.current) {
      const editor = diffEditorRef.current;
      const modifiedEditor = editor.getModifiedEditor();
      const originalEditor = editor.getOriginalEditor();

      if (originalEditor.getValue() !== original) {
        originalEditor.setValue(original);
      }

      if (modifiedEditor.getValue() !== modified) {
        modifiedEditor.setValue(modified);
      }
    }
  }, [original, modified]);

  return (
    <Card className="h-full flex flex-col gap-2">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center">
          <span className="w-6 h-6 bg-blue-600 rounded-full flex items-center justify-center text-white text-sm">📝</span>
          Template Diff Viewer
        </CardTitle>
      </CardHeader>

      {/* Status & Error Display */}
      {(error || !wasmLoaded) && (
        <div className="px-6 pb-3">
          <div className={`flex items-center gap-3 p-3 rounded-lg ${error ? 'bg-red-50 border border-red-200 text-red-700' : 'bg-blue-50 border border-blue-200 text-blue-700'}`}>
            <div className={`w-6 h-6 rounded-full flex items-center justify-center text-white text-sm ${
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
          <div className="flex-1 min-h-[400px]">
            <DiffEditor
              height="100%"
              original={localOriginal}
              modified={modified}
              language="html"
              theme="vs-dark"
              options={{
                readOnly,
                minimap: { enabled: false },
                fontSize: 14,
                lineNumbers: 'on',
                scrollBeyondLastLine: false,
                automaticLayout: true,
                wordWrap: 'on',
                renderSideBySide: true,
                originalEditable: !readOnly,
                renderOverviewRuler: true,
                enableSplitViewResizing: true,
              }}
              onMount={(editor) => {
                diffEditorRef.current = editor;
                console.log('DiffEditor mounted');

                // Set up change listener for the original editor when not read-only
                if (!readOnly) {
                  const originalEditor = editor.getOriginalEditor();
                  console.log('Setting up original editor change listener');

                  // Listen for content changes
                  const disposable = originalEditor.onDidChangeModelContent(() => {
                    const value = originalEditor.getValue();
                    console.log('Original editor content changed:', value);
                    handleEditorChange(value);
                  });

                  // Cleanup on unmount
                  return () => {
                    console.log('Cleaning up editor listeners');
                    disposable.dispose();
                  };
                }
              }}
            />
          </div>
          <div className="flex justify-between items-center px-4 py-2 bg-muted text-sm text-muted-foreground border-t">
            <div className="text-xs">
              {readOnly ? 'Read-only Diff View' : 'Editable Template with Diff View'}
            </div>
            <div className="text-xs flex items-center gap-4">
              <div className="flex items-center gap-1">
                <div className="w-3 h-3 bg-red-500 rounded"></div>
                <span>Original Template</span>
              </div>
              <div className="flex items-center gap-1">
                <div className="w-3 h-3 bg-green-500 rounded"></div>
                <span>Rendered Output</span>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}