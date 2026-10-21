// Author: HyperparameterSI
#import <Cocoa/Cocoa.h>
#include "shell.h"
#include "include/cef_version.h"

namespace {
NSColor* Color(unsigned value) {
  return [NSColor colorWithSRGBRed:((value >> 16) & 255) / 255.0
                           green:((value >> 8) & 255) / 255.0
                            blue:(value & 255) / 255.0 alpha:1];
}
NSString* Text(const std::string& text) { return [NSString stringWithUTF8String:text.c_str()] ?: @""; }
NSString* Navigation(NSString* input) {
  NSString* text = [input stringByTrimmingCharactersInSet:NSCharacterSet.whitespaceAndNewlineCharacterSet];
  if (!text.length) return nil;
  if ([text hasPrefix:@"https://"] || [text hasPrefix:@"http://"]) return text;
  BOOL spaces = [text rangeOfCharacterFromSet:NSCharacterSet.whitespaceAndNewlineCharacterSet].location != NSNotFound;
  if (!spaces && ([text containsString:@"."] || [text hasPrefix:@"localhost"])) {
    return [([text hasPrefix:@"localhost"] || [text hasPrefix:@"127.0.0.1"]) ? @"http://" : @"https://" stringByAppendingString:text];
  }
  NSCharacterSet* allowed = [NSCharacterSet characterSetWithCharactersInString:@"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"];
  return [@"https://www.google.com/search?q=" stringByAppendingString:[text stringByAddingPercentEncodingWithAllowedCharacters:allowed]];
}
}

@class GoletBrowserShell;
@interface GoletShellView : NSView
@property(nonatomic, weak) GoletBrowserShell* owner;
@end

@interface GoletBrowserShell : NSObject <NSWindowDelegate, NSTextFieldDelegate>
@property(nonatomic, strong) NSWindow* window;
@property(nonatomic, strong) GoletShellView* root;
@property(nonatomic, strong) NSView* pages;
@property(nonatomic, strong) NSView* addressBackground;
@property(nonatomic, strong) NSTextField* address;
@property(nonatomic, strong) NSTextField* engine;
@property(nonatomic, strong) NSTextField* brand;
@property(nonatomic, strong) NSImageView* searchIcon;
@property(nonatomic, strong) NSButton* back;
@property(nonatomic, strong) NSButton* forward;
@property(nonatomic, strong) NSButton* reload;
@property(nonatomic, strong) NSButton* add;
@property(nonatomic, strong) NSButton* home;
@property(nonatomic, strong) NSButton* assistant;
@property(nonatomic, strong) NSButton* fullscreen;
@property(nonatomic, strong) NSMutableArray<NSNumber*>* order;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSView*>* views;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSView*>* tabs;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSButton*>* labels;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSButton*>* closes;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSDictionary*>* states;
@property(nonatomic) NSInteger selected;
@property(nonatomic) BOOL closing;
- (void)layout;
- (void)updateAddress;
- (void)focusAddress:(id)sender;
- (void)newTab:(id)sender;
- (void)closeActive:(id)sender;
- (void)navigate:(id)sender;
@end

@implementation GoletShellView
- (BOOL)isFlipped { return YES; }
- (BOOL)mouseDownCanMoveWindow { return YES; }
- (void)layout { [super layout]; [self.owner layout]; }
- (void)drawRect:(NSRect)dirtyRect {
  [Color(0xfbfcfa) setFill]; NSRectFill(self.bounds);
  [Color(0xf0f3ee) setFill]; NSRectFill(NSMakeRect(0, 0, self.bounds.size.width, 50));
  [Color(0xe5e9df) setFill]; NSRectFill(NSMakeRect(0, 113, self.bounds.size.width, 1));
}
@end

@implementation GoletBrowserShell
- (NSButton*)button:(NSString*)symbol label:(NSString*)label action:(SEL)action {
  NSButton* button = [NSButton buttonWithImage:[NSImage imageWithSystemSymbolName:symbol accessibilityDescription:label] target:self action:action];
  button.bordered = NO;
  button.contentTintColor = Color(0x87967c);
  button.toolTip = label;
  [button setAccessibilityLabel:label];
  [self.root addSubview:button];
  return button;
}
- (instancetype)init {
  if (!(self = [super init])) return nil;
  self.order = [NSMutableArray array]; self.views = [NSMutableDictionary dictionary];
  self.tabs = [NSMutableDictionary dictionary]; self.labels = [NSMutableDictionary dictionary];
  self.closes = [NSMutableDictionary dictionary]; self.states = [NSMutableDictionary dictionary];
  self.window = [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, 1200, 820)
     styleMask:NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable | NSWindowStyleMaskFullSizeContentView
     backing:NSBackingStoreBuffered defer:NO];
  self.window.releasedWhenClosed = NO;
  self.window.delegate = self;
  self.window.title = @"Golet Browser";
  self.window.titleVisibility = NSWindowTitleHidden;
  self.window.titlebarAppearsTransparent = YES;
  self.window.minSize = NSMakeSize(720, 500);
  self.window.backgroundColor = Color(0xfbfcfa);
  self.window.appearance = [NSAppearance appearanceNamed:NSAppearanceNameAqua];
  self.root = [[GoletShellView alloc] initWithFrame:NSMakeRect(0, 0, 1200, 820)];
  self.root.owner = self;
  self.window.contentView = self.root;
  self.pages = [[NSView alloc] initWithFrame:NSZeroRect];
  [self.root addSubview:self.pages];
  self.back = [self button:@"arrow.left" label:@"后退" action:@selector(goBack:)];
  self.forward = [self button:@"arrow.right" label:@"前进" action:@selector(goForward:)];
  self.reload = [self button:@"arrow.clockwise" label:@"刷新" action:@selector(refresh:)];
  self.assistant = [self button:@"sparkles" label:@"AI 助手与扩展" action:@selector(assistant:)];
  self.home = [self button:@"house" label:@"首页" action:@selector(goHome:)];
  self.fullscreen = [self button:@"arrow.up.left.and.arrow.down.right" label:@"全屏" action:@selector(fullscreen:)];
  self.add = [self button:@"plus" label:@"新建标签页" action:@selector(newTab:)];
  self.brand = [NSTextField labelWithString:@"golet browser"];
  self.brand.font = [NSFont systemFontOfSize:12 weight:NSFontWeightRegular]; self.brand.textColor = Color(0x9aa493);
  [self.root addSubview:self.brand];
  self.addressBackground = [[NSView alloc] initWithFrame:NSZeroRect];
  self.addressBackground.wantsLayer = YES; self.addressBackground.layer.cornerRadius = 20;
  self.addressBackground.layer.backgroundColor = Color(0xf0f3ec).CGColor;
  self.addressBackground.layer.borderWidth = 1; self.addressBackground.layer.borderColor = Color(0xe3e9db).CGColor;
  [self.root addSubview:self.addressBackground];
  self.searchIcon = [[NSImageView alloc] initWithFrame:NSZeroRect];
  self.searchIcon.image = [NSImage imageWithSystemSymbolName:@"magnifyingglass" accessibilityDescription:nil];
  self.searchIcon.contentTintColor = Color(0x98a58b); [self.addressBackground addSubview:self.searchIcon];
  self.address = [[NSTextField alloc] initWithFrame:NSZeroRect];
  self.address.bordered = NO; self.address.drawsBackground = NO; self.address.focusRingType = NSFocusRingTypeNone;
  self.address.font = [NSFont systemFontOfSize:14]; self.address.textColor = Color(0x46513e);
  self.address.placeholderString = @"搜索或输入网址"; self.address.delegate = self;
  self.address.target = self; self.address.action = @selector(navigate:);
  [self.address setAccessibilityLabel:@"地址栏"];
  [self.addressBackground addSubview:self.address];
  self.engine = [NSTextField labelWithString:@"Google"];
  self.engine.font = [NSFont systemFontOfSize:12]; self.engine.textColor = Color(0x95a286);
  [self.addressBackground addSubview:self.engine];
  [self installMenu];
  [self.window center]; [self layout]; [self.window makeKeyAndOrderFront:nil];
  [NSApp activateIgnoringOtherApps:YES];
  return self;
}
- (void)installMenu {
  NSMenu* main = [[NSMenu alloc] init];
  NSMenuItem* appItem = [[NSMenuItem alloc] init]; [main addItem:appItem];
  NSMenu* appMenu = [[NSMenu alloc] initWithTitle:@"Golet Browser"];
  [appMenu addItemWithTitle:@"退出 Golet Browser" action:@selector(terminate:) keyEquivalent:@"q"];
  appItem.submenu = appMenu;
  NSMenuItem* fileItem = [[NSMenuItem alloc] init]; [main addItem:fileItem];
  NSMenu* fileMenu = [[NSMenu alloc] initWithTitle:@"文件"];
  NSArray* titles = @[@"新建标签页", @"关闭标签页", @"地址栏", @"刷新"];
  NSArray* keys = @[@"t", @"w", @"l", @"r"];
  SEL actions[] = {@selector(newTab:), @selector(closeActive:), @selector(focusAddress:), @selector(refresh:)};
  for (NSUInteger i = 0; i < titles.count; ++i) {
    NSMenuItem* item = [fileMenu addItemWithTitle:titles[i] action:actions[i] keyEquivalent:keys[i]]; item.target = self;
  }
  fileItem.submenu = fileMenu;
  NSMenuItem* editItem = [[NSMenuItem alloc] init]; [main addItem:editItem];
  NSMenu* edit = [[NSMenu alloc] initWithTitle:@"编辑"];
  [edit addItemWithTitle:@"撤销" action:@selector(undo:) keyEquivalent:@"z"];
  [edit addItemWithTitle:@"剪切" action:@selector(cut:) keyEquivalent:@"x"];
  [edit addItemWithTitle:@"复制" action:@selector(copy:) keyEquivalent:@"c"];
  [edit addItemWithTitle:@"粘贴" action:@selector(paste:) keyEquivalent:@"v"];
  [edit addItemWithTitle:@"全选" action:@selector(selectAll:) keyEquivalent:@"a"];
  editItem.submenu = edit; NSApp.mainMenu = main;
}
- (void)layout {
  CGFloat width = self.root.bounds.size.width, height = self.root.bounds.size.height;
  self.pages.frame = NSMakeRect(0, 114, width, MAX(0, height - 114));
  for (NSView* view in self.views.allValues) view.frame = self.pages.bounds;
  self.back.frame = NSMakeRect(18, 65, 30, 32); self.forward.frame = NSMakeRect(54, 65, 30, 32);
  self.reload.frame = NSMakeRect(90, 65, 30, 32);
  self.addressBackground.frame = NSMakeRect(138, 61, width - 292, 42);
  CGFloat addressWidth = self.addressBackground.bounds.size.width;
  self.searchIcon.frame = NSMakeRect(16, 12, 18, 18);
  self.address.frame = NSMakeRect(46, 11, addressWidth - 118, 21);
  self.engine.frame = NSMakeRect(addressWidth - 61, 13, 49, 18);
  self.assistant.frame = NSMakeRect(width - 138, 65, 32, 32);
  self.home.frame = NSMakeRect(width - 96, 65, 32, 32); self.fullscreen.frame = NSMakeRect(width - 54, 65, 32, 32);
  self.brand.frame = NSMakeRect(width - 120, 18, 104, 18);
  CGFloat tabWidth = MIN(186, (width - 312) / MAX((NSUInteger)1, self.order.count));
  CGFloat x = 108;
  for (NSNumber* key in self.order) {
    NSView* tab = self.tabs[key]; tab.frame = NSMakeRect(x, 8, tabWidth, 42);
    self.labels[key].frame = NSMakeRect(12, 8, MAX(20, tabWidth - 46), 27);
    self.closes[key].frame = NSMakeRect(tabWidth - 29, 12, 20, 20);
    x += tabWidth + 3;
  }
  self.add.frame = NSMakeRect(MIN(x + 8, width - 162), 12, 28, 30);
}
- (void)updateAddress {
  NSDictionary* state = self.states[@(self.selected)];
  NSString* url = state[@"url"] ?: @"";
  BOOL home = [url hasPrefix:Text(ShellHomeURL())] || [url isEqualToString:@"about:blank"];
  // Do not overwrite text while the user is editing the location field.
  if (!self.address.currentEditor) self.address.stringValue = home ? @"" : url;
  self.back.enabled = [state[@"back"] boolValue]; self.forward.enabled = [state[@"forward"] boolValue];
  BOOL loading = [state[@"loading"] boolValue];
  self.reload.image = [NSImage imageWithSystemSymbolName:loading ? @"xmark" : @"arrow.clockwise" accessibilityDescription:loading ? @"停止加载" : @"刷新"];
  self.window.title = home ? @"Golet Browser" : [NSString stringWithFormat:@"%@ · Golet Browser", state[@"title"] ?: @"网页"];
}
- (BOOL)windowShouldClose:(NSWindow*)sender { if (self.closing) return YES; BrowserUIAction("close"); return NO; }
- (void)windowDidResize:(NSNotification*)notification { [self layout]; }
- (void)newTab:(id)sender { BrowserUIAction("newtab"); }
- (void)closeActive:(id)sender { BrowserUIAction("close-tab", (int)self.selected); }
- (void)chooseTab:(NSButton*)sender { BrowserUIAction("select-tab", (int)sender.tag); }
- (void)closeTab:(NSButton*)sender { BrowserUIAction("close-tab", (int)sender.tag); }
- (void)goBack:(id)sender { BrowserUIAction("back"); }
- (void)goForward:(id)sender { BrowserUIAction("forward"); }
- (void)refresh:(id)sender { BrowserUIAction([self.states[@(self.selected)][@"loading"] boolValue] ? "stop" : "reload"); }
- (void)assistant:(id)sender { BrowserUIAction("assistant"); }
- (void)goHome:(id)sender { BrowserUIAction("home"); }
- (void)fullscreen:(id)sender { [self.window toggleFullScreen:sender]; }
- (void)focusAddress:(id)sender { [self.window makeFirstResponder:self.address]; [self.address selectText:nil]; }
- (void)navigate:(id)sender {
  NSString* url = Navigation(self.address.stringValue);
  if (url) { [self.window makeFirstResponder:self.pages]; BrowserUIAction("navigate", 0, url.UTF8String); }
}
@end

namespace { GoletBrowserShell* shell; }
void CreateBrowserShell() { shell = [[GoletBrowserShell alloc] init]; }
void ConfigureShellChild(CefWindowInfo& info) {
  NSRect bounds = shell.pages.bounds;
  info.SetAsChild((__bridge CefWindowHandle)shell.pages, CefRect(0, 0, bounds.size.width, bounds.size.height));
  info.runtime_style = CEF_RUNTIME_STYLE_ALLOY;
}
void ShellAttachBrowser(int id, CefWindowHandle handle) {
  NSNumber* key = @(id); NSView* view = (__bridge NSView*)handle;
  shell.views[key] = view; view.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
  [shell.order addObject:key];
  NSView* tab = [[NSView alloc] initWithFrame:NSZeroRect]; tab.wantsLayer = YES; tab.layer.cornerRadius = 11;
  shell.tabs[key] = tab; [shell.root addSubview:tab];
  NSButton* label = [NSButton buttonWithTitle:@"新标签页" target:shell action:@selector(chooseTab:)];
  label.bordered = NO; label.alignment = NSTextAlignmentLeft; label.font = [NSFont systemFontOfSize:13];
  label.image = [NSImage imageWithSystemSymbolName:@"globe" accessibilityDescription:nil]; label.imagePosition = NSImageLeft;
  label.contentTintColor = Color(0x56624d); label.tag = id;
  label.cell.lineBreakMode = NSLineBreakByTruncatingTail;
  shell.labels[key] = label; [tab addSubview:label];
  NSButton* close = [NSButton buttonWithImage:[NSImage imageWithSystemSymbolName:@"xmark" accessibilityDescription:@"关闭标签页"] target:shell action:@selector(closeTab:)];
  close.bordered = NO; close.contentTintColor = Color(0x809075); close.tag = id;
  shell.closes[key] = close; [tab addSubview:close];
  [shell layout];
}
void ShellDetachBrowser(int id) {
  [shell.views[@(id)] removeFromSuperview];
  [shell.views removeObjectForKey:@(id)];
}
void ShellRemoveBrowser(int id) {
  NSNumber* key = @(id); [shell.tabs[key] removeFromSuperview];
  [shell.tabs removeObjectForKey:key]; [shell.views removeObjectForKey:key];
  [shell.labels removeObjectForKey:key]; [shell.closes removeObjectForKey:key];
  [shell.states removeObjectForKey:key]; [shell.order removeObject:key]; [shell layout];
}
void ShellSelectBrowser(int id) {
  shell.selected = id;
  for (NSNumber* key in shell.order) {
    BOOL selected = key.intValue == id;
    shell.views[key].hidden = !selected;
    shell.tabs[key].layer.backgroundColor = Color(selected ? 0xfbfcfa : 0xf0f3ee).CGColor;
  }
  [shell updateAddress];
}
void ShellUpdateBrowser(int id, const std::string& title, const std::string& url, bool loading, bool back, bool forward) {
  NSNumber* key = @(id);
  shell.states[key] = @{@"title":Text(title), @"url":Text(url), @"loading":@(loading), @"back":@(back), @"forward":@(forward)};
  shell.labels[key].title = title.empty() ? @"新标签页" : Text(title);
  shell.labels[key].toolTip = Text(title);
  if (shell.selected == id) [shell updateAddress];
}
void ShellDestroy() { shell.closing = YES; [shell.window close]; shell = nil; }
std::string ShellHomeURL() {
  NSURL* url = [[[NSBundle mainBundle] resourceURL] URLByAppendingPathComponent:@"newtab/index.html"];
  NSString* version = [NSString stringWithFormat:@"?version=%d.%d.%d.%d", CHROME_VERSION_MAJOR, CHROME_VERSION_MINOR, CHROME_VERSION_BUILD, CHROME_VERSION_PATCH];
  return [[url.absoluteString stringByAppendingString:version] UTF8String];
}
