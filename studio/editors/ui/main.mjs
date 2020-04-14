import "/views/ui/vnd/mithril-2.0.4.min.js?0";
import "/views/ui/vnd/qtalk/qmux.js?0";
import "/views/ui/vnd/qtalk/qrpc.js?0";

import * as hotweb from '/views/.hotweb/client.mjs'
import * as app from '/views/ui/lib/app.js';

function wrap(cb) {
    return { view: () => m(cb()) };
}

hotweb.watchCSS();
hotweb.watchHTML();
hotweb.refresh(() => m.redraw())
m.mount(document.body, wrap(() => app.App));

