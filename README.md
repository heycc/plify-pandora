# Plify Pandora 🧰

> A collection of interesting and useful tools built by via Vibe Coding

## 📦 Tools Collection

---

## 🎯 Go Template Live - Write Go Templates with Confidence

**Even if you're not a Go developer**

An interactive web-based playground for Go templates that instantly shows what variables you need, previews your output, and validates syntax - all in your browser. Perfect for SREs, DevOps engineers, and anyone working with Go template-based systems (Helm, Confd, Consul Template) who just wants to get their template working.

**Location:** `go-template-live/`

### ✨ Key Features

- **🔍 Auto-Extract Variables** - Instantly see all variables referenced in your template, or show you where the syntax fail.
- **🎁 Smart Default Detection** - Automatically extracts default values from custom functions like `{{getv "key" "default"}}
- **📊 Live Diff View** - See exactly what changes as you type, with side-by-side comparison using Monaco Editor.
- **⚡ All Real-Time** - Extracting and rendering your template instantly with WebAssembly-powered processing.
- **🔗 URL Sharing** - Share your template content via URL for quick team review and collaboration.
- **🛠️ Confd-Style Functions** - Built-in support for `getv`, `exists`, `get`, `json`, `jsonArray` .etc that up to 20+ custom funcitons.

### 🚀 Why Use This?

#### Traditional Workflow:
- ❌ Write template
- ❌ Deploy to test environment
- ❌ Find error: "variable 'UserName' not found"
- ❌ Fix and redeploy
- ❌ Another error: "expected string, got bool"
- ❌ Fix and redeploy again...


#### With Go Template Live:
- ✅ Paste template → See variables needed
- ✅ Fill in test values → Preview output
- ✅ See diff in real-time → Iterate quickly
- ✅ Copy working template → Deploy with confidence