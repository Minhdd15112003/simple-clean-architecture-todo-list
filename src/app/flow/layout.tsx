import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import "../../css/flow.css";
import AppSidebar from "@/components/flow/AppSidebar";
export default function FlowLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div>
      <SidebarProvider>
        <AppSidebar />
        <SidebarTrigger />
        <main>{children}</main>
      </SidebarProvider>
    </div>
  );
}
