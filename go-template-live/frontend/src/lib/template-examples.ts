// Template examples - stored separately to avoid JavaScript template literal parsing issues

export const examples = {
  Basic: `
Hello {{.Name}}, your email is {{.Email}}!

Feel free to try the powerful golang template engine.

This app are able to:
- extract the variables from the template
- render the template with the variables
- show the diff
all in real-time.
`.trim(),

Control_Flow: `
{{if .Authenticated}}
  Welcome {{.User}}!
{{else}}
  Please login to continue.
{{end}}
`.trim(),

Default_Value: `
# Configuration
server={{getv "server" "localhost"}}
port={{getv "port" "8080"}}
debug={{getv "debug" "false"}}
{{if (exists "database")}}
  database={{.database}}
{{end}}

{{- /* Note: \`getv\` is a custom function, it takes a key name and default value as input. \`exists\` is a custom function.*/}}
`.trim(),

  Loop: `
{{range $key, $val := json "MapData"}}
Key: {{$key}}, Value: {{$val}}
{{- end }}

---
{{- /* Note: \`json\` is a custom function, it takes a JSON string as input and returns a map. */}}
{{- /* Try to input type {"Apple": "260", "Microsoft": "510"} at right panel and see the result. */}}
`.trim(),

  Nested_Fields: `
{{with json "Config"}}
    Server: {{.host}}
    Port: {{.port}}
    {{- if .ssl}}
        HTTPS Enabled
    {{end}}
{{end}}

{{- /* Note: \`json\` is a custom function, it takes a JSON string as input and returns a map. */}}
{{- /* Try to input type {"host": "127.0.0.1", "port":80} at right panel and see the result. */}}
`.trim(),

  Email: `
Dear {{getv "name" "Customer"}},
{{if (exists "order_id")}}
Thank you for your order #{{.order_id}}!
{{end}}
Your items:
{{- range $item := jsonArray "items"}}
- {{$item.name}}: \${{$item.price -}}
{{end}}

Total: \${{getv "total" "N/A"}}


Best regards,
{{getv "company" "Our Company"}}

{{- /* Note: \`jsonArray\`, \`exists\`, \`getv\` are custom functions */}}
{{- /* Try to input items like [{"name": "Apple", "price": 1.99}, {"name": "Banana", "price": 0.99}] at right panel and see the result. */}}
`.trim(),

  Format_Control: `
Comments will not rendered.{{/* this is a comment */}}

{{\`"a raw output with quotes"\`}}
  `.trim(),
};