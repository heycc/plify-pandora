'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import Editor from '@monaco-editor/react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface TemplateDiffViewerProps {
  original: string;
  modified: string;
  onOriginalChange?: (value: string) => void;
  readOnly?: boolean;
  error?: string | null;
  wasmLoaded?: boolean;
}

interface DiffRange {
  startLine: number;
  startColumn: number;
  endLine: number;
  endColumn: number;
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
  const monacoRef = useRef<any>(null);
  const [decorations, setDecorations] = useState<{ original: any[], modified: any[] }>({ original: [], modified: [] });
  const originalDecorationsRef = useRef<string[]>([]);
  const modifiedDecorationsRef = useRef<string[]>([]);

  // LCS-based diff algorithm to find and highlight changes
  const calculateDiffRanges = useCallback((original: string, modified: string) => {
    const originalRanges: DiffRange[] = [];
    const modifiedRanges: DiffRange[] = [];

    // Use LCS (Longest Common Subsequence) for better diff
    const lcs = computeLCS(original, modified);
    
    // Build ranges from LCS
    let origIndex = 0;
    let modIndex = 0;
    let lcsIndex = 0;
    let diffStartOrig = -1;
    let diffStartMod = -1;

    while (origIndex < original.length || modIndex < modified.length) {
      const lcsChar = lcsIndex < lcs.length ? lcs[lcsIndex] : null;
      const origChar = origIndex < original.length ? original[origIndex] : null;
      const modChar = modIndex < modified.length ? modified[modIndex] : null;

      if (origChar === lcsChar && modChar === lcsChar) {
        // Both match LCS, save any pending diff
        if (diffStartOrig !== -1) {
          const pos = indexToLineColumn(original, diffStartOrig);
          const endPos = indexToLineColumn(original, origIndex);
          originalRanges.push({
            startLine: pos.line,
            startColumn: pos.column,
            endLine: endPos.line,
            endColumn: endPos.column
          });
          diffStartOrig = -1;
        }
        if (diffStartMod !== -1) {
          const pos = indexToLineColumn(modified, diffStartMod);
          const endPos = indexToLineColumn(modified, modIndex);
          modifiedRanges.push({
            startLine: pos.line,
            startColumn: pos.column,
            endLine: endPos.line,
            endColumn: endPos.column
          });
          diffStartMod = -1;
        }
        origIndex++;
        modIndex++;
        lcsIndex++;
      } else if (origChar !== lcsChar && origChar !== null) {
        // Character in original that's different
        if (diffStartOrig === -1) diffStartOrig = origIndex;
        origIndex++;
      } else if (modChar !== lcsChar && modChar !== null) {
        // Character in modified that's different
        if (diffStartMod === -1) diffStartMod = modIndex;
        modIndex++;
      } else {
        break;
      }
    }

    // Save any final diff
    if (diffStartOrig !== -1) {
      const pos = indexToLineColumn(original, diffStartOrig);
      const endPos = indexToLineColumn(original, origIndex);
      originalRanges.push({
        startLine: pos.line,
        startColumn: pos.column,
        endLine: endPos.line,
        endColumn: endPos.column
      });
    }
    if (diffStartMod !== -1) {
      const pos = indexToLineColumn(modified, diffStartMod);
      const endPos = indexToLineColumn(modified, modIndex);
      modifiedRanges.push({
        startLine: pos.line,
        startColumn: pos.column,
        endLine: endPos.line,
        endColumn: endPos.column
      });
    }

    return { originalRanges, modifiedRanges };
  }, []);

  // Compute Longest Common Subsequence
  const computeLCS = (str1: string, str2: string): string => {
    const m = str1.length;
    const n = str2.length;
    const dp: number[][] = Array(m + 1).fill(0).map(() => Array(n + 1).fill(0));

    for (let i = 1; i <= m; i++) {
      for (let j = 1; j <= n; j++) {
        if (str1[i - 1] === str2[j - 1]) {
          dp[i][j] = dp[i - 1][j - 1] + 1;
        } else {
          dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1]);
        }
      }
    }

    // Backtrack to find LCS
    let lcs = '';
    let i = m, j = n;
    while (i > 0 && j > 0) {
      if (str1[i - 1] === str2[j - 1]) {
        lcs = str1[i - 1] + lcs;
        i--;
        j--;
      } else if (dp[i - 1][j] > dp[i][j - 1]) {
        i--;
      } else {
        j--;
      }
    }

    return lcs;
  };

  // Helper function to convert string index to line/column
  const indexToLineColumn = (text: string, index: number): { line: number; column: number } => {
    const beforeText = text.substring(0, index);
    const lines = beforeText.split('\n');
    return {
      line: lines.length,
      column: lines[lines.length - 1].length + 1
    };
  };

  // Update decorations when content changes
  useEffect(() => {
    if (!originalEditorRef.current || !modifiedEditorRef.current || !monacoRef.current) return;

    const { originalRanges, modifiedRanges } = calculateDiffRanges(original, modified);
    const monaco = monacoRef.current;

    // Create decorations for original editor (red background for removed/changed parts)
    const originalDecorations = originalRanges.map(range => ({
      range: new monaco.Range(range.startLine, range.startColumn, range.endLine, range.endColumn),
      options: {
        isWholeLine: false,
        className: 'diff-removed-line',
        inlineClassName: 'diff-removed-inline',
        backgroundColor: 'rgba(255, 100, 100, 0.2)',
        minimap: {
          color: 'rgba(255, 100, 100, 0.6)',
          position: monaco.editor.MinimapPosition.Inline
        }
      }
    }));

    // Create decorations for modified editor (green background for added/changed parts)
    const modifiedDecorations = modifiedRanges.map(range => ({
      range: new monaco.Range(range.startLine, range.startColumn, range.endLine, range.endColumn),
      options: {
        isWholeLine: false,
        className: 'diff-added-line',
        inlineClassName: 'diff-added-inline',
        backgroundColor: 'rgba(100, 255, 100, 0.2)',
        minimap: {
          color: 'rgba(100, 255, 100, 0.6)',
          position: monaco.editor.MinimapPosition.Inline
        }
      }
    }));

    setDecorations({ original: originalDecorations, modified: modifiedDecorations });
  }, [original, modified, calculateDiffRanges]);

  // Apply decorations
  useEffect(() => {
    if (originalEditorRef.current) {
      const editor = originalEditorRef.current;
      originalDecorationsRef.current = editor.deltaDecorations(
        originalDecorationsRef.current,
        decorations.original
      );
    }
    if (modifiedEditorRef.current) {
      const editor = modifiedEditorRef.current;
      modifiedDecorationsRef.current = editor.deltaDecorations(
        modifiedDecorationsRef.current,
        decorations.modified
      );
    }
  }, [decorations]);

  const handleOriginalMount = useCallback((editor: any, monaco: any) => {
    originalEditorRef.current = editor;
    monacoRef.current = monaco;
  }, []);

  const handleModifiedMount = useCallback((editor: any, monaco: any) => {
    modifiedEditorRef.current = editor;
    if (!monacoRef.current) {
      monacoRef.current = monaco;
    }
  }, []);

  const handleOriginalChange = useCallback((value: string | undefined) => {
    if (onOriginalChange && value !== undefined && !readOnly) {
      onOriginalChange(value);
    }
  }, [onOriginalChange, readOnly]);

  // Sync scrolling between editors
  useEffect(() => {
    if (!originalEditorRef.current || !modifiedEditorRef.current) return;

    const originalEditor = originalEditorRef.current;
    const modifiedEditor = modifiedEditorRef.current;

    const originalScrollListener = originalEditor.onDidScrollChange((e: any) => {
      modifiedEditor.setScrollPosition({ scrollTop: e.scrollTop });
    });

    const modifiedScrollListener = modifiedEditor.onDidScrollChange((e: any) => {
      originalEditor.setScrollPosition({ scrollTop: e.scrollTop });
    });

    return () => {
      originalScrollListener.dispose();
      modifiedScrollListener.dispose();
    };
  }, []);

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
          {/* Labels */}
          <div className="flex border-b bg-muted">
            <div className="flex-1 px-4 py-2 text-sm font-medium border-r">
              Template (Editable)
            </div>
            <div className="flex-1 px-4 py-2 text-sm font-medium">
              Rendered Result
            </div>
          </div>

          {/* Side-by-side editors */}
          <div className="flex-1 min-h-[400px] flex">
            <div className="flex-1 border-r">
              <Editor
                height="100%"
                value={original}
                language="html"
                theme="vs-dark"
                options={{
                  readOnly: readOnly,
                  minimap: { enabled: false },
                  fontSize: 14,
                  lineNumbers: 'on',
                  scrollBeyondLastLine: false,
                  automaticLayout: true,
                  wordWrap: 'on',
                  folding: false,
                  glyphMargin: false,
                  lineDecorationsWidth: 10,
                  lineNumbersMinChars: 3,
                }}
                onChange={handleOriginalChange}
                onMount={handleOriginalMount}
              />
            </div>
            <div className="flex-1">
              <Editor
                height="100%"
                value={modified}
                language="html"
                theme="vs-dark"
                options={{
                  readOnly: true,
                  minimap: { enabled: false },
                  fontSize: 14,
                  lineNumbers: 'on',
                  scrollBeyondLastLine: false,
                  automaticLayout: true,
                  wordWrap: 'on',
                  folding: false,
                  glyphMargin: false,
                  lineDecorationsWidth: 10,
                  lineNumbersMinChars: 3,
                }}
                onMount={handleModifiedMount}
              />
            </div>
          </div>
          
          <div className="flex justify-between items-center px-4 py-2 bg-muted text-sm text-muted-foreground border-t">
            <div className="text-xs">
              {readOnly ? 'Read-only View' : '✏️ Edit template on left to see live rendering on right'}
            </div>
            <div className="text-xs flex items-center gap-4">
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded" style={{ backgroundColor: 'rgba(255, 100, 100, 0.3)' }}></div>
                <span>Template Syntax</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded" style={{ backgroundColor: 'rgba(100, 255, 100, 0.3)' }}></div>
                <span>Rendered Value</span>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}