// Template examples - stored separately to avoid JavaScript template literal parsing issues

export const examples = {
  Basic: 'Hello {{.Name}}, your email is {{.Email}}!',
  Control_Flow: '{{if .Authenticated}}\n    Welcome {{.User.Name}}!\n    {{range .Items}}\n        - {{.Name}}: {{.Price}}\n    {{end}}\n{{else}}\n    Please login to continue.\n{{end}}'
};

export const defaultExamples = {
  Nested_Fields: '{{with .Config}}\n    Server: {{getv .server "localhost"}}\n    Port: {{getv .port "8080"}}\n    {{if (getv .ssl false)}}\n        HTTPS Enabled\n    {{end}}\n{{end}}',
  Custom_Functions: 'Server: {{getv "server" "localhost"}}:{{getv "port" "8080"}}\nDatabase: {{getv "database" "mysql"}}\nDebug: {{getv "debug" "false"}}'
};

export const renderExamples = {
  profile: '<div class="profile">\n    <h2>{{getv "name" "Anonymous"}}</h2>\n    <p>Email: {{getv "email" "no-email@example.com"}}</p>\n    <p>Role: {{getv "role" "User"}}</p>\n    {{if (exists "department")}}\n    <p>Department: {{.department}}</p>\n    {{end}}\n</div>',
  config: '# Configuration\nserver={{getv "server" "localhost"}}\nport={{getv "port" "8080"}}\ndebug={{getv "debug" "false"}}\n{{if (exists "database")}}\ndatabase={{.database}}\n{{end}}',
  email: 'Hello {{getv "name" "Customer"}},\n\n{{if (exists "order_id")}}\nThank you for your order #{{.order_id}}!\n{{end}}\n\nYour items:\n{{range .items}}\n- {{.name}}: ${{.price}}\n{{end}}\n\nTotal: ${{getv "total" "0.00"}}\n\nBest regards,\n{{getv "company" "Our Company"}}'
};