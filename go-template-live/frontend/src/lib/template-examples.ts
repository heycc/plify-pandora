// Template examples - stored separately to avoid JavaScript template literal parsing issues

export const examples = {
  simple: 'Hello {{.Name}}, your email is {{.Email}}!',
  complex: '<!DOCTYPE html>\n<html>\n<head><title>{{.Title}}</title></head>\n<body>\n    <h1>Welcome {{.User.Name}}</h1>\n    <p>Email: {{.User.Email}}</p>\n    {{if .User.Active}}\n        <p>Status: Active</p>\n    {{end}}\n</body>\n</html>',
  functions: 'Server: {{getv "server" "localhost"}}:{{getv "port" "8080"}}\nDatabase: {{getv "database" "mysql"}}\nDebug: {{getv "debug" "false"}}',
  control: '{{if .Authenticated}}\n    Welcome {{.User.Name}}!\n    {{range .Items}}\n        - {{.Name}}: {{.Price}}\n    {{end}}\n{{else}}\n    Please login to continue.\n{{end}}'
};

export const defaultExamples = {
  basic: 'Welcome {{getv "username" "guest"}}!\nEmail: {{getv "email" "no-reply@example.com"}}\nTheme: {{getv "theme" "light"}}',
  nested: '{{with .Config}}\n    Server: {{getv .server "localhost"}}\n    Port: {{getv .port "8080"}}\n    {{if (getv .ssl false)}}\n        HTTPS Enabled\n    {{end}}\n{{end}}'
};

export const renderExamples = {
  profile: '<div class="profile">\n    <h2>{{getv "name" "Anonymous"}}</h2>\n    <p>Email: {{getv "email" "no-email@example.com"}}</p>\n    <p>Role: {{getv "role" "User"}}</p>\n    {{if (exists "department")}}\n    <p>Department: {{.department}}</p>\n    {{end}}\n</div>',
  config: '# Configuration\nserver={{getv "server" "localhost"}}\nport={{getv "port" "8080"}}\ndebug={{getv "debug" "false"}}\n{{if (exists "database")}}\ndatabase={{.database}}\n{{end}}',
  email: 'Hello {{getv "name" "Customer"}},\n\n{{if (exists "order_id")}}\nThank you for your order #{{.order_id}}!\n{{end}}\n\nYour items:\n{{range .items}}\n- {{.name}}: ${{.price}}\n{{end}}\n\nTotal: ${{getv "total" "0.00"}}\n\nBest regards,\n{{getv "company" "Our Company"}}'
};