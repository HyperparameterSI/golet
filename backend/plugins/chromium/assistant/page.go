// Author: HyperparameterSI
package assistant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HyperparameterSI/golet/backend/plugins/chromium"
	"net/url"
)

// Page is an explicit, bounded snapshot. Form values and password fields are excluded.
type Page struct {
	URL      string    `json:"url"`
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	Elements []Element `json:"elements"`
}
type Element struct {
	ID    int    `json:"id"`
	Role  string `json:"role"`
	Label string `json:"label"`
}

// Action is a proposed operation. Execute only after user review of this exact action.
type Action struct {
	Kind    string `json:"kind"`
	Element int    `json:"element"`
	Value   string `json:"value,omitempty"`
}

const snapshotScript = `(() => {
 const nodes=[...document.querySelectorAll('a,button,input,textarea,select,[role="button"]')].filter(e=>e.getClientRects().length && !e.disabled && e.type!=='password' && e.type!=='hidden').slice(0,100);
 const elements=nodes.map((e,id)=>({id,role:e.tagName.toLowerCase(),label:(e.getAttribute('aria-label')||e.innerText||e.getAttribute('placeholder')||e.getAttribute('name')||'').slice(0,120)}));
 const clone=document.body.cloneNode(true); clone.querySelectorAll('script,style,input,textarea,select,[contenteditable]').forEach(e=>e.remove());
 return {url:location.href,title:document.title,text:clone.innerText?.slice(0,16000)||clone.textContent.slice(0,16000),elements};
})()`

func Snapshot(ctx context.Context, b *chromium.Browser, tab int) (Page, error) {
	var page Page
	err := b.Evaluate(ctx, tab, snapshotScript, &page)
	if err == nil {
		u, e := url.Parse(page.URL)
		if e != nil || (u.Scheme != "https" && u.Scheme != "http") {
			err = errors.New("only HTTP(S) pages can be shared")
		}
	}
	return page, err
}

// Apply rechecks the exact page URL and element label before performing a fixed
// action. It never evaluates model-provided code or accepts a model selector.
func Apply(ctx context.Context, b *chromium.Browser, tab int, page Page, action Action) error {
	if action.Kind != "click" && action.Kind != "fill" && action.Kind != "scroll" {
		return errors.New("unsupported action")
	}
	label := ""
	if action.Kind != "scroll" {
		if action.Element < 0 || action.Element >= len(page.Elements) {
			return errors.New("invalid element")
		}
		label = page.Elements[action.Element].Label
	}
	payload, _ := json.Marshal(map[string]any{"url": page.URL, "action": action, "label": label})
	script := fmt.Sprintf(`(() => {const p=%s;if(location.href!==p.url)throw Error('Page changed');
 const nodes=[...document.querySelectorAll('a,button,input,textarea,select,[role="button"]')].filter(e=>e.getClientRects().length && !e.disabled && e.type!=='password' && e.type!=='hidden').slice(0,100);
 if(p.action.kind==='scroll'){window.scrollBy(0,innerHeight*0.8);return true}
 const e=nodes[p.action.element];if(!e)throw Error('Element missing');
 const label=(e.getAttribute('aria-label')||e.innerText||e.getAttribute('placeholder')||e.getAttribute('name')||'').slice(0,120);if(label!==p.label)throw Error('Element changed');
 if(p.action.kind==='click')e.click();else{if(!['INPUT','TEXTAREA'].includes(e.tagName)||['password','file','hidden'].includes(e.type))throw Error('Cannot fill');const proto=e.tagName==='INPUT'?HTMLInputElement.prototype:HTMLTextAreaElement.prototype;Object.getOwnPropertyDescriptor(proto,'value').set.call(e,p.action.value);e.dispatchEvent(new Event('input',{bubbles:true}));e.dispatchEvent(new Event('change',{bubbles:true}));}return true;})()`, payload)
	return b.Evaluate(ctx, tab, script, nil)
}
