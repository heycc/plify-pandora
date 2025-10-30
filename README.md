# Plify Trove 🧰

> A collection of interesting and useful tools built by **Vibe Coding**

## 📦 Tools Collection

---

## 🎯 Go Template Live - Write Go Templates with Confidence

**Even if you're not a Go developer**

An interactive web-based playground for Go templates that instantly shows what variables you need, previews your output, and validates syntax - all in your browser. Perfect for SREs, DevOps engineers, and anyone working with Go template-based systems (Helm, Confd, Consul Template) who just wants to get their template working.

**Location:** `go-template-live/`

### ✨ Key Features

- **🔍 Auto-Extract Variables** - Stop guessing what data your template needs. Instantly see all variables referenced in your template.
- **📊 Live Diff View** - See exactly what changes as you type, with side-by-side comparison using Monaco Editor.
- **⚡ Real-Time Rendering** - Preview your template output instantly with WebAssembly-powered processing.
- **🎁 Smart Default Detection** - Automatically extracts default values from custom functions like `{{getv "key" "default"}}`.
- **🔗 URL Sharing** - Share templates via URL for quick team review and collaboration.
- **🛠️ Confd-Style Functions** - Built-in support for `getv`, `exists`, `get`, `json`, `jsonArray` functions.

### 🎯 Perfect For

| Use Case | Description |
|----------|-------------|
| 📦 **Helm Chart Templates** | Test `values.yaml` rendering before deployment |
| ⚙️ **Confd Configurations** | Preview config file output with different variables |
| 📧 **Email Templates** | See variables and formatting in real-time |
| 🔧 **Custom Config Systems** | Validate template syntax before deployment |
| 📚 **Learning Go Templates** | Experiment and learn by seeing immediate results |

### 🚀 Why Use This?

#### Traditional Workflow:
❌ Write template
❌ Deploy to test environment
❌ Find error: "variable 'UserName' not found"
❌ Fix and redeploy
❌ Another error: "expected string, got bool"
❌ Fix and redeploy again...


#### With Go Template Live:
✅ Paste template → See variables needed
✅ Fill in test values → Preview output
✅ See diff in real-time → Iterate quickly
✅ Copy working template → Deploy with confidence


✨ Live Preview: Updates as you type
🔄 Diff View: Seeitecture for adding new template functions
- **Monaco Diff Editor** - Professional code editing experience
- **Next.js + React** - Modern, responsive frontend
- **Zero Server Dependency** - All processing happens client-side