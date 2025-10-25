This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Hosting Under a Subdirectory

This project supports hosting under a subdirectory (e.g., `/template-parser`). To configure this:

1. Set the `NEXT_PUBLIC_BASE_PATH` environment variable before building:
   ```bash
   export NEXT_PUBLIC_BASE_PATH=/template-parser
   npm run build
   ```

2. Or create a `.env.local` file:
   ```
   NEXT_PUBLIC_BASE_PATH=/template-parser
   ```

3. The base path is automatically applied to:
   - All static assets (CSS, JS, images)
   - WASM resources (`wasm_exec.js`, `main.wasm`)
   - Internal navigation links

4. When deploying with nginx, make sure your nginx configuration matches the base path. See `nginx.conf.example` for reference.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.
