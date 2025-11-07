const fs = require('fs');
const path = require('path');

// Define source and destination paths
const monacoDir = path.join(__dirname, '../node_modules/monaco-editor');
const publicDir = path.join(__dirname, '../public');
const monacoPublicDir = path.join(publicDir, 'monaco-editor');

// Function to recursively copy directory
function copyDir(src, dest) {
  // Create destination directory if it doesn't exist
  if (!fs.existsSync(dest)) {
    fs.mkdirSync(dest, { recursive: true });
  }

  // Read all files/folders in source directory
  const entries = fs.readdirSync(src, { withFileTypes: true });

  for (let entry of entries) {
    const srcPath = path.join(src, entry.name);
    const destPath = path.join(dest, entry.name);

    if (entry.isDirectory()) {
      // Recursively copy subdirectory
      copyDir(srcPath, destPath);
    } else {
      // Copy file
      fs.copyFileSync(srcPath, destPath);
    }
  }
}

try {
  console.log('Copying Monaco Editor files to public directory...');
  
  // Copy the min directory (contains all the necessary files)
  const monacoMinDir = path.join(monacoDir, 'min');
  if (fs.existsSync(monacoMinDir)) {
    copyDir(monacoMinDir, monacoPublicDir);
    console.log('✅ Monaco Editor files copied successfully!');
  } else {
    console.error('❌ Monaco Editor min directory not found at:', monacoMinDir);
    process.exit(1);
  }
} catch (error) {
  console.error('❌ Error copying Monaco Editor files:', error);
  process.exit(1);
}

