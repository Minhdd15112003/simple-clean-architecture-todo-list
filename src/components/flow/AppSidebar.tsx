"use client";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupAction,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar";
import { Calendar, ChevronRight, Home, Inbox, Plus, Search, Settings } from "lucide-react";
import { Button } from "../ui/button";
import { Collapsible } from "@radix-ui/react-collapsible";
import { CollapsibleContent, CollapsibleTrigger } from "../ui/collapsible";
// Menu items.
const items = [
  {
    label: "loop",
  },
  {
    label: "pause",
  },
  {
    label: "sleep",
  },
  {
    label: "toast",
  },
  {
    label: "unlockScreen",
  },
  {
    label: "pressHome",
  },
  {
    label: "pressBack",
  },
  {
    label: "pressSwitch",
  },
  {
    label: "getClipboard",
  },
  {
    label: "setClipboard",
  },
  {
    label: "appStart",
  },
  {
    label: "appStop",
  },
  {
    label: "appClear",
  },
  {
    label: "swipe",
  },
  {
    label: "keyCode",
  },
  {
    label: "keyText",
  },
  {
    label: "resolution",
  },
  {
    label: "findElements",
  },
];
export default function AppSidebar() {
  return (
    <Sidebar>
      <SidebarHeader>
        <SidebarGroupLabel>Execute</SidebarGroupLabel>
        <Button
          variant='secondary'
          size='sm'
          className='w-full mb-2'
          onClick={() => {
            const event = new CustomEvent("add-node", { detail: "start" });
            window.dispatchEvent(event); // có thể lắng nghe sự kiện này ở component cha
          }}
        >
          Start
        </Button>
        <Button
          variant='secondary'
          size='sm'
          className='w-full mb-2'
          onClick={() => {
            const event = new CustomEvent("download", { detail: "download" });
            window.dispatchEvent(event); // có thể lắng nghe sự kiện này ở component cha
          }}
        >
          Save flow
        </Button>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Nodes</SidebarGroupLabel>
          <Collapsible defaultOpen className='group/collapsible'>
            {/* <SidebarMenuItem> */}
            <CollapsibleTrigger asChild>
              <SidebarMenuButton tooltip='Node Actions'>
                <span>Node Actions</span>
                <ChevronRight className='ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90' />
              </SidebarMenuButton>
            </CollapsibleTrigger>
            <CollapsibleContent>
              {/* <SidebarMenuSub> */}
              {items.map((item) => (
                <Button
                  key={item.label}
                  variant='secondary'
                  size='sm'
                  className='w-full mb-2'
                  onClick={() => {
                    const event = new CustomEvent("add-node", { detail: item.label });
                    window.dispatchEvent(event); // có thể lắng nghe sự kiện này ở component cha
                  }}
                >
                  {/* <Plus className='mr-2' /> */}
                  {item.label}
                </Button>
              ))}
              {/* </SidebarMenuSub> */}
            </CollapsibleContent>
            {/* </SidebarMenuItem> */}
          </Collapsible>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
}
