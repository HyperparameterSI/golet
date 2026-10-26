// Author: HyperparameterSI
#import <Cocoa/Cocoa.h>
#include "shell.h"
#include "include/cef_version.h"
#include "include/cef_command_line.h"

namespace {
NSDictionary* theme;
NSColor* Color(unsigned value) {
  NSString* role = nil;
  switch (value) {
    case 0xfbfcfa: role = @"background"; break;
    case 0xf0f3ee: role = @"surface"; break;
    case 0xf0f3ec: role = @"muted"; break;
    case 0xe5e9df: case 0xe3e9db: role = @"border"; break;
    case 0x46513e: case 0x56624d: role = @"foreground"; break;
    case 0x87967c: role = @"primary"; break;
    case 0x98a58b: case 0x95a286: case 0x809075: role = @"mutedForeground"; break;
  }
  NSString* hex = role ? theme[role] : nil;
  if (hex) {
    unsigned long long rgba = 0;
    [[NSScanner scannerWithString:hex] scanHexLongLong:&rgba];
    return [NSColor colorWithSRGBRed:((rgba >> 24) & 255) / 255.0
      green:((rgba >> 16) & 255) / 255.0 blue:((rgba >> 8) & 255) / 255.0 alpha:(rgba & 255) / 255.0];
  }
  return [NSColor colorWithSRGBRed:((value >> 16) & 255) / 255.0
                           green:((value >> 8) & 255) / 255.0
                            blue:(value & 255) / 255.0 alpha:1];
}
NSString* Text(const std::string& text) { return [NSString stringWithUTF8String:text.c_str()] ?: @""; }
NSString* ProductName() {
  auto name = CefCommandLine::GetGlobalCommandLine()->GetSwitchValue("golet-app-name");
  if (!name.empty()) return Text(name.ToString());
  NSString* bundled = [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleDisplayName"];
  return bundled.length ? bundled : @"Golet Browser";
}
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
@property(nonatomic, strong) NSImageView* searchIcon;
@property(nonatomic, strong) NSButton* back;
@property(nonatomic, strong) NSButton* forward;
@property(nonatomic, strong) NSButton* reload;
@property(nonatomic, strong) NSButton* add;
@property(nonatomic, strong) NSButton* home;
@property(nonatomic, strong) NSButton* assistant;
@property(nonatomic, strong) NSButton* computerUse;
@property(nonatomic, strong) NSButton* fullscreen;
@property(nonatomic, strong) NSMutableArray<NSNumber*>* order;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSView*>* views;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSView*>* tabs;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSButton*>* labels;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSButton*>* closes;
@property(nonatomic, strong) NSMutableDictionary<NSNumber*, NSDictionary*>* states;
@property(nonatomic) NSInteger selected;
@property(nonatomic) BOOL closing;
@property(nonatomic) BOOL appMode;
- (void)applyTheme;
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
  if (self.owner.appMode) return;
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
  self.appMode = CefCommandLine::GetGlobalCommandLine()->HasSwitch("golet-app-mode");
  self.order = [NSMutableArray array]; self.views = [NSMutableDictionary dictionary];
  self.tabs = [NSMutableDictionary dictionary]; self.labels = [NSMutableDictionary dictionary];
  self.closes = [NSMutableDictionary dictionary]; self.states = [NSMutableDictionary dictionary];
  self.window = [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, 1200, 820)
     styleMask:NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable | (self.appMode ? 0 : NSWindowStyleMaskFullSizeContentView)
     backing:NSBackingStoreBuffered defer:NO];
  self.window.releasedWhenClosed = NO;
  self.window.delegate = self;
  self.window.title = ProductName();
  self.window.titleVisibility = self.appMode ? NSWindowTitleVisible : NSWindowTitleHidden;
  self.window.titlebarAppearsTransparent = !self.appMode;
  self.window.minSize = NSMakeSize(720, 500);
  self.window.backgroundColor = Color(0xfbfcfa);
  self.window.appearance = [NSAppearance appearanceNamed:NSAppearanceNameAqua];
  self.root = [[GoletShellView alloc] initWithFrame:NSMakeRect(0, 0, 1200, 820)];
  self.root.owner = self;
  self.window.contentView = self.root;
  self.pages = [[NSView alloc] initWithFrame:NSZeroRect];
  [self.root addSubview:self.pages];
  if (!self.appMode) {
    self.back = [self button:@"arrow.left" label:@"后退" action:@selector(goBack:)];
    self.forward = [self button:@"arrow.right" label:@"前进" action:@selector(goForward:)];
    self.reload = [self button:@"arrow.clockwise" label:@"刷新" action:@selector(refresh:)];
    self.assistant = [self button:@"sparkles" label:@"AI 助手与扩展" action:@selector(assistant:)];
    if (CefCommandLine::GetGlobalCommandLine()->HasSwitch("golet-computer-use")) {
      self.computerUse = [NSButton buttonWithTitle:@"AI 操作：关闭" target:self action:@selector(toggleComputerUse:)];
      [self.computerUse setButtonType:NSButtonTypePushOnPushOff];
      self.computerUse.state = NSControlStateValueOff;
      [self.computerUse setAccessibilityLabel:@"允许 AI 操作此浏览器"];
      self.computerUse.toolTip = @"仅授权当前浏览器窗口；再次点击立即关闭 AI 操作";
      [self.root addSubview:self.computerUse];
    }
    self.home = [self button:@"house" label:@"首页" action:@selector(goHome:)];
    self.fullscreen = [self button:@"arrow.up.left.and.arrow.down.right" label:@"全屏" action:@selector(fullscreen:)];
    self.add = [self button:@"plus" label:@"新建标签页" action:@selector(newTab:)];
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
  }
  [self installMenu];
  [self applyTheme]; [self.window center]; [self layout]; [self.window makeKeyAndOrderFront:nil];
  [NSApp activateIgnoringOtherApps:YES];
  return self;
}
- (void)applyTheme {
  if (!theme) return;
  self.window.appearance = [NSAppearance appearanceNamed:[theme[@"dark"] boolValue] ? NSAppearanceNameDarkAqua : NSAppearanceNameAqua];
  self.window.backgroundColor = Color(0xfbfcfa);
  CGFloat size = [theme[@"fontSize"] doubleValue];
  NSFont* font = [NSFont fontWithName:theme[@"font"] size:size] ?: [NSFont systemFontOfSize:size];
  self.address.font = font;
  self.address.textColor = Color(0x46513e);
  self.addressBackground.layer.cornerRadius = MIN(21, [theme[@"radius"] doubleValue]);
  self.addressBackground.layer.backgroundColor = Color(0xf0f3ec).CGColor;
  self.addressBackground.layer.borderColor = Color(0xe3e9db).CGColor;
  self.engine.textColor = Color(0x95a286);
  self.searchIcon.contentTintColor = Color(0x98a58b);
  for (NSView* view in self.root.subviews) {
    if ([view isKindOfClass:NSButton.class]) ((NSButton*)view).contentTintColor = Color(0x87967c);
  }
  for (NSNumber* key in self.order) {
    self.tabs[key].layer.cornerRadius = MIN(21, [theme[@"radius"] doubleValue]);
    self.tabs[key].layer.backgroundColor = Color(key.intValue == self.selected ? 0xfbfcfa : 0xf0f3ee).CGColor;
    self.labels[key].font = font;
    self.labels[key].contentTintColor = Color(0x56624d);
    self.closes[key].contentTintColor = Color(0x809075);
  }
  self.root.needsDisplay = YES;
}
- (void)installMenu {
  NSMenu* main = [[NSMenu alloc] init];
  NSMenuItem* appItem = [[NSMenuItem alloc] init]; [main addItem:appItem];
  NSMenu* appMenu = [[NSMenu alloc] initWithTitle:ProductName()];
  [appMenu addItemWithTitle:[@"退出 " stringByAppendingString:ProductName()] action:@selector(terminate:) keyEquivalent:@"q"];
  appItem.submenu = appMenu;
  NSMenuItem* fileItem = [[NSMenuItem alloc] init]; [main addItem:fileItem];
  NSMenu* fileMenu = [[NSMenu alloc] initWithTitle:@"文件"];
  NSArray* titles = @[@"新建标签页", @"关闭标签页", @"地址栏", @"刷新"];
  NSArray* keys = @[@"t", @"w", @"l", @"r"];
  SEL actions[] = {@selector(newTab:), @selector(closeActive:), @selector(focusAddress:), @selector(refresh:)};
  for (NSUInteger i = 0; i < titles.count; ++i) {
    if (self.appMode && (i == 0 || i == 2)) continue;
    NSMenuItem* item = [fileMenu addItemWithTitle:(self.appMode && i == 1) ? @"关闭窗口" : titles[i] action:actions[i] keyEquivalent:keys[i]]; item.target = self;
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
  if (self.appMode) {
    self.pages.frame = self.root.bounds;
    for (NSView* view in self.views.allValues) view.frame = self.pages.bounds;
    return;
  }
  self.pages.frame = NSMakeRect(0, 114, width, MAX(0, height - 114));
  for (NSView* view in self.views.allValues) view.frame = self.pages.bounds;
  self.back.frame = NSMakeRect(18, 65, 30, 32); self.forward.frame = NSMakeRect(54, 65, 30, 32);
  self.reload.frame = NSMakeRect(90, 65, 30, 32);
  self.addressBackground.frame = NSMakeRect(138, 61, width - (self.computerUse ? 408 : 292), 42);
  CGFloat addressWidth = self.addressBackground.bounds.size.width;
  self.searchIcon.frame = NSMakeRect(16, 12, 18, 18);
  self.address.frame = NSMakeRect(46, 11, addressWidth - 118, 21);
  self.engine.frame = NSMakeRect(addressWidth - 61, 13, 49, 18);
  self.computerUse.frame = NSMakeRect(width - 260, 65, 112, 32);
  self.assistant.frame = NSMakeRect(width - 138, 65, 32, 32);
  self.home.frame = NSMakeRect(width - 96, 65, 32, 32); self.fullscreen.frame = NSMakeRect(width - 54, 65, 32, 32);
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
  if (self.appMode) {
    NSString* title = state[@"title"];
    self.window.title = title.length ? title : ProductName();
    return;
  }
  NSString* url = state[@"url"] ?: @"";
  BOOL home = [url hasPrefix:Text(ShellHomeURL())] || [url isEqualToString:@"about:blank"];
  // Do not overwrite text while the user is editing the location field.
  if (!self.address.currentEditor) self.address.stringValue = home ? @"" : url;
  self.back.enabled = [state[@"back"] boolValue]; self.forward.enabled = [state[@"forward"] boolValue];
  BOOL loading = [state[@"loading"] boolValue];
  self.reload.image = [NSImage imageWithSystemSymbolName:loading ? @"xmark" : @"arrow.clockwise" accessibilityDescription:loading ? @"停止加载" : @"刷新"];
  self.window.title = home ? ProductName() : [NSString stringWithFormat:@"%@ · %@", state[@"title"] ?: @"网页", ProductName()];
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
- (void)toggleComputerUse:(id)sender {
  self.computerUse.title = self.computerUse.state == NSControlStateValueOn ? @"AI 操作：开启" : @"AI 操作：关闭";
}
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
bool ShellSetTheme(const std::string& json) {
  NSData* data = [Text(json) dataUsingEncoding:NSUTF8StringEncoding];
  id parsed = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
  if (![parsed isKindOfClass:NSDictionary.class]) return false;
  for (NSString* key in @[@"background", @"surface", @"muted", @"foreground", @"mutedForeground", @"border", @"primary"]) {
    id color = parsed[key];
    if (![color isKindOfClass:NSString.class] || [color length] != 8 ||
        [color rangeOfCharacterFromSet:[[NSCharacterSet characterSetWithCharactersInString:@"0123456789abcdefABCDEF"] invertedSet]].location != NSNotFound) return false;
  }
  for (NSString* key in @[@"radius", @"fontSize", @"dark"]) if (![parsed[key] isKindOfClass:NSNumber.class]) return false;
  if (![parsed[@"font"] isKindOfClass:NSString.class]) return false;
  double radius = [parsed[@"radius"] doubleValue], size = [parsed[@"fontSize"] doubleValue];
  if (!(radius >= 0 && radius <= 64 && size >= 10 && size <= 20)) return false;
  theme = parsed;
  [shell applyTheme];
  return true;
}
void CreateBrowserShell() {
  auto json = CefCommandLine::GetGlobalCommandLine()->GetSwitchValue("golet-theme");
  if (!json.empty()) ShellSetTheme(json.ToString());
  shell = [[GoletBrowserShell alloc] init];
}
void ConfigureShellChild(CefWindowInfo& info) {
  NSRect bounds = shell.pages.bounds;
  info.SetAsChild((__bridge CefWindowHandle)shell.pages, CefRect(0, 0, bounds.size.width, bounds.size.height));
  info.runtime_style = CEF_RUNTIME_STYLE_ALLOY;
}
void ShellAttachBrowser(int id, CefWindowHandle handle) {
  NSNumber* key = @(id); NSView* view = (__bridge NSView*)handle;
  shell.views[key] = view; view.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
  [shell.order addObject:key];
  if (shell.appMode) { [shell layout]; return; }
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
  [shell applyTheme];
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
  NSURLComponents* components = [NSURLComponents componentsWithURL:url resolvingAgainstBaseURL:YES];
  components.queryItems = @[[NSURLQueryItem queryItemWithName:@"version" value:[version substringFromIndex:9]],
                            [NSURLQueryItem queryItemWithName:@"appName" value:ProductName()]];
  return components.URL.absoluteString.UTF8String;
}

// Author: HyperparameterSI. CEF supplies an NSView, while the driver requires
// its owning WindowServer window number. Never return the renderer process ID.
int ShellNativeWindowID(CefWindowHandle handle) {
  NSView* view = (__bridge NSView*)handle;
  NSWindow* window = view.window;
  return window && window.visible ? (int)window.windowNumber : 0;
}
bool ShellComputerUseEnabled() {
  return shell && !shell.appMode && shell.computerUse && shell.computerUse.state == NSControlStateValueOn;
}
