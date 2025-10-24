// Template examples - stored separately to avoid JavaScript template literal parsing issues

export const examples = {
  Basic: `
Hello {{.Name}}, your email is {{.Email}}!

Feel free to try the powerful golang template engine.

This app are able to extract the variables from the template, render the template with the variables, show the diff, all in real-time.
`.trim(),

Control_Flow: `
{{if .Authenticated}}
  Welcome {{.User}}!
{{else}}
  Please login to continue.
{{end}}
`.trim(),

  Loop: `
Note: \`json\` is a custom function, it takes a JSON string as input and returns a map.
Try to input type {"Apple": "260", "Microsoft": "510"} at right panel and see the result.

{{range $key, $val := json "MapData"}}
Key: {{$key}}, Value: {{$val}}
{{- end }}
`.trim(),

  Nested_Fields: `
{{with .Config}}
    Server: {{.host}}
    Port: {{.port}}
    {{if .ssl}}
        HTTPS Enabled
    {{end}}
{{end}}
`.trim(),


  Format_Control: `
Comments will not rendered.{{/* this is a comment */}}

{{\`"a raw output with quotes"\`}}
  `.trim(),


  Custom_Functions: 'Server: {{getv "server" "localhost"}}:{{getv "port" "8080"}}\nDatabase: {{getv "database" "mysql"}}\nDebug: {{getv "debug" "false"}}',
  Config: '# Configuration\nserver={{getv "server" "localhost"}}\nport={{getv "port" "8080"}}\ndebug={{getv "debug" "false"}}\n{{if (exists "database")}}\ndatabase={{.database}}\n{{end}}',
  Email: 'Hello {{getv "name" "Customer"}},\n\n{{if (exists "order_id")}}\nThank you for your order #{{.order_id}}!\n{{end}}\n\nYour items:\n{{range .items}}\n- {{.name}}: ${{.price}}\n{{end}}\n\nTotal: ${{getv "total" "0.00"}}\n\nBest regards,\n{{getv "company" "Our Company"}}'
};