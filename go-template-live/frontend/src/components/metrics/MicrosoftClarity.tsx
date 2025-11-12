'use client'

import Script from "next/script"

const MicrosoftClarity = () => {
    // Check if Microsoft Clarity ID is configured
    if (!process.env.NEXT_PUBLIC_MICROSOFT_CLARITY) {
        return (
            <>
                {/* Microsoft Clarity not configured. Set NEXT_PUBLIC_MICROSOFT_CLARITY environment variable to enable. */}
            </>
        )
    }

    return (
        <Script
            id="microsoft-clarity-init"
            strategy="afterInteractive"
            dangerouslySetInnerHTML={{
                __html: `
                (function(c,l,a,r,i,t,y){
                    c[a]=c[a]||function(){(c[a].q=c[a].q||[]).push(arguments)};
                    t=l.createElement(r);t.async=1;t.src="https://www.clarity.ms/tag/"+i;
                    y=l.getElementsByTagName(r)[0];y.parentNode.insertBefore(t,y);
                })(window, document, "clarity", "script", "${process.env.NEXT_PUBLIC_MICROSOFT_CLARITY}");
                `,
            }}
        />
    )
}

export default MicrosoftClarity