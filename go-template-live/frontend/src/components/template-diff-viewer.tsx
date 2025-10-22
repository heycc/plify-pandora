'use client';

import { useState, useEffect, useRef } from 'react';
import { Editor } from '@monaco-editor/react';
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
  const originalEditorRef = useRef<any>(null);
  const modifiedEditorRef = useRef<any>(null);

  // Update editor content directly when props change
  useEffect(() => {
    if (originalEditorRef.current && originalEditorRef.current.getValue() !== original) {
      originalEditorRef.current.setValue(original);
    }
  }, [original]);

  useEffect(() => {
    if (modifiedEditorRef.current && modifiedEditorRef.current.getValue() !== modified) {
      modifiedEditorRef.current.setValue(modified);
    }
  }, [modified]);

  return (
    <Card className="h-full flex flex-col gap-2">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center">
          <span className="w-6 h-6 bg-blue-600 rounded-full flex items-center justify-center text-white text-sm">📝</span>
          Template Editor & Preview
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
          <div className="flex-1 min-h-[400px] grid grid-cols-2">
            {/* Left: Original Template */}
            <div className="border-r">
              <div className="h-8 flex items-center px-4 bg-muted text-sm font-medium border-b">
                <div className="w-3 h-3 bg-red-500 rounded mr-2"></div>
                Original Template
              </div>
              <Editor
                height="calc(100% - 2rem)"
                defaultLanguage="html"
                theme="vs-dark"
                options={{
                  readOnly,
                  minimap: { enabled: false },
                  fontSize: 14,
                  lineNumbers: 'on',
                  scrollBeyondLastLine: false,
                  automaticLayout: true,
                  wordWrap: 'on',
                }}
                value={original}
                onChange={(value) => onOriginalChange?.(value || '')}
                onMount={(editor) => {
                  originalEditorRef.current = editor;
                }}
              />
            </div>

            {/* Right: Rendered Output */}
            <div>
              <div className="h-8 flex items-center px-4 bg-muted text-sm font-medium border-b">
                <div className="w-3 h-3 bg-green-500 rounded mr-2"></div>
                Rendered Output
              </div>
              <Editor
                height="calc(100% - 2rem)"
                defaultLanguage="html"
                theme="vs-dark"
                options={{
                  readOnly: true,
                  minimap: { enabled: false },
                  fontSize: 14,
                  lineNumbers: 'on',
                  scrollBeyondLastLine: false,
                  automaticLayout: true,
                  wordWrap: 'on',
                }}
                value={modified}
                onMount={(editor) => {
                  modifiedEditorRef.current = editor;
                }}
              />
            </div>
          </div>
          <div className="flex justify-between items-center px-4 py-2 bg-muted text-sm text-muted-foreground border-t">
            <div className="text-xs">
              {readOnly ? 'Read-only Preview' : 'Editable Template'}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}