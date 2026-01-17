"use client";

import {
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  Background,
  Connection,
  Controls,
  type Edge,
  type EdgeTypes,
  IsValidConnection,
  type Node,
  type OnConnect,
  type OnEdgesChange,
  type OnNodesChange,
  ReactFlow,
} from "@xyflow/react";
import { useCallback, useEffect, useState } from "react";
import "@xyflow/react/dist/style.css";
import NewNode from "@/common/nodeHandler";
import nodeTypes from "@/components/flow/customNodes";
import edgeTypes from "@/components/flow/customEdges";

export default function FlowPage() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  //get nodes from event
  useEffect(() => {
    const handleAddNode = (event: any) => {
      const nodeLabel = event.detail;
      setNodes((nds) => {
        const newNode: Node = new NewNode(nds.length, {
          label: nodeLabel,
        }).CreateNode();
        return nds.concat(newNode);
      });
    };

    window.addEventListener("add-node", handleAddNode);

    return () => {
      window.removeEventListener("add-node", handleAddNode);
    };
  }, []);

  useEffect(() => {
    const handleDownload = (event: any) => {
      const flowData = {
        nodes: nodes,
        edges: edges,
      };
      const dataStr =
        "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(flowData, null, 2));
      const downloadAnchorNode = document.createElement("a");
      downloadAnchorNode.setAttribute("href", dataStr);
      downloadAnchorNode.setAttribute("download", "flow_data.json");
      document.body.appendChild(downloadAnchorNode); // required for firefox
      downloadAnchorNode.click();
      downloadAnchorNode.remove();
    };

    window.addEventListener("download", handleDownload);

    return () => {
      window.removeEventListener("download", handleDownload);
    };
  }, [nodes, edges]);

  const onNodesChange: OnNodesChange = useCallback(
    (changes) => setNodes((nodesSnapshot) => applyNodeChanges(changes, nodesSnapshot)),
    [],
  );
  const onEdgesChange: OnEdgesChange = useCallback(
    (changes) => setEdges((edgesSnapshot) => applyEdgeChanges(changes, edgesSnapshot)),
    [],
  );
  const onConnect: OnConnect = useCallback((params) => {
    const newEdge = { ...params, type: "custom" };

    setEdges((eds) => addEdge(newEdge, eds));
  }, []);

  const onNodesDelete = useCallback((deleted: Node[]) => {
    console.log("Deleted nodes:", deleted);
  }, []);

  const onEdgesDelete = useCallback((deleted: Edge[]) => {
    console.log("Deleted edges:", deleted);
  }, []);

  const isValidConnection = useCallback((connection: Edge | Connection): boolean => {
    const isHorizontal =
      (connection.sourceHandle === "right" && connection.targetHandle === "left") ||
      (connection.sourceHandle === "left" && connection.targetHandle === "right");
    const isVertical =
      (connection.sourceHandle === "bottom" && connection.targetHandle === "top") ||
      (connection.sourceHandle === "top" && connection.targetHandle === "bottom");

    const isValid = isHorizontal || isVertical;

    if (!isValid) {
      alert("Invalid connection: Right must connect to Left, and Top must connect to Bottom.");
    }

    return isValid;
  }, []);

  return (
    <div style={{ width: "100vw", height: "100vh" }}>
      <ReactFlow
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        isValidConnection={isValidConnection}
        onNodesDelete={onNodesDelete}
        onEdgesDelete={onEdgesDelete}
        deleteKeyCode={["Backspace", "Delete"]}
        colorMode='dark'
      >
        <Background />
        <Controls />
      </ReactFlow>
    </div>
  );
}
