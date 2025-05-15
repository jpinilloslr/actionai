import * as PanelMenu from "resource:///org/gnome/shell/ui/panelMenu.js";
import * as PopupMenu from "resource:///org/gnome/shell/ui/popupMenu.js";
import { Extension } from "resource:///org/gnome/shell/extensions/extension.js";
import St from "gi://St";
import GObject from "gi://GObject";
import * as Main from "resource:///org/gnome/shell/ui/main.js";

const Indicator = GObject.registerClass(
  class Indicator extends PanelMenu.Button {
    _init() {
      super._init(0.0, _("ActionAI Indicator"));

      this.add_child(
        new St.Icon({
          icon_name: "view-mirror-symbolic",
          style_class: "system-status-icon",
        }),
      );

      const item = new PopupMenu.PopupMenuItem(_("Show Notification"));
      item.connect("activate", () => {
        Main.notify(_("This is atest"));
      });
      (this.menu as any).addMenuItem(item);
    }
  },
);

export default class ActionAiExtension extends Extension {
  _indicator: PanelMenu.Button | null = null;

  enable() {
    this._indicator = new Indicator(0.0, this.metadata.name, false);

    const icon = new St.Icon({
      icon_name: "face-laugh-symbolic",
      style_class: "system-status-icon",
    });
    this._indicator.add_child(icon);

    Main.panel.addToStatusArea(this.uuid, this._indicator);
  }

  disable() {
    this._indicator?.destroy();
    this._indicator = null;
  }
}
