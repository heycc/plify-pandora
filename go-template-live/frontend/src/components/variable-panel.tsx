'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';

interface VariableInfo {
  name: string;
  defaultValue?: string;
}

interface VariablePanelProps {
  variables: VariableInfo[];
  values: Record<string, string>;
  onVariablesChange: (variables: Record<string, string>) => void;
}

export function VariablePanel({ variables, values, onVariablesChange }: VariablePanelProps) {
  const handleVariableChange = (variableName: string, value: string) => {
    const newValues = {
      ...values,
      [variableName]: value
    };
    onVariablesChange(newValues);
  };

  if (variables.length === 0) {
    return (
      <Card className="h-full gap-2 py-4 bg-white/70 backdrop-blur-sm">
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <span className="text-2xl">🎯</span>
            Variables
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center text-muted-foreground py-8">
            <div className="text-4xl mb-2">🔍</div>
            <p>No variables found in template</p>
            <p className="text-sm mt-1">Start typing a template to extract variables</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="h-full flex flex-col gap-2 py-4 bg-gray-100 backdrop-blur-sm">
      <CardHeader className="">
        <CardTitle className="text-lg flex items-center gap-2">
          <span className="text-2xl">🎯</span>
          Variables ({variables.length})
        </CardTitle>
      </CardHeader>
      <CardContent className="flex-1 p-0">
        <div className="h-full overflow-y-auto">
          <div className="space-y-4 px-4">
            {variables.map((variable, index) => (
              <div key={`${variable.name}-${index}`} className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor={`variable-${variable.name}-${index}`} className="font-mono text-sm">
                    {variable.name}
                  </Label>
                  {variable.defaultValue && (
                    <Badge variant="secondary" className="text-xs">
                      Default: "{variable.defaultValue}"
                    </Badge>
                  )}
                </div>
                <Input
                  id={`variable-${variable.name}-${index}`}
                  value={values[variable.name] || ''}
                  onChange={(e) => handleVariableChange(variable.name, e.target.value)}
                  placeholder={variable.defaultValue ? `Default: ${variable.defaultValue}` : 'Enter value...'}
                  className="font-mono text-sm py-0"
                />
              </div>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}