import {
  BaseEdge,
  EdgeLabelRenderer,
  getBezierPath,
  useStore,
  type EdgeProps,
} from "@xyflow/react";
import { FC } from "react";

export const CustomEdge: FC<EdgeProps> = ({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  label,
  style,
  markerEnd,
}) => {
  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  // Example of how to get selection status
  const isSelected = useStore((s) => s.edges.find((e) => e.id === id)?.selected);
  // console.log({
  //   id,
  //   sourceX,
  //   sourceY,
  //   targetX,
  //   targetY,
  //   sourcePosition,
  //   targetPosition,
  //   label,
  //   style,
  //   markerEnd,
  //   isSelected,
  // });
  if (targetPosition === "left" && sourcePosition === "right") {
    label = "Success";
  } else if (targetPosition === "top" && sourcePosition === "bottom") {
    label = "Failure";
  }
  return (
    <>
      <BaseEdge
        id={id}
        path={edgePath}
        markerEnd={markerEnd}
        style={{
          ...style,
          stroke: isSelected ? "#6366F1" : style?.stroke,
          strokeWidth: isSelected ? 2 : 1,
        }}
      />

      <EdgeLabelRenderer>
        <div
          style={{
            position: "absolute",
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            pointerEvents: "all",
            ...(label == "Success"
              ? { backgroundColor: "green" }
              : label == "Failure"
                ? { backgroundColor: "red" }
                : { backgroundColor: "gray" }),
          }}
          className='nodrag nopan rounded-md bg-gray-800 px-2 py-1 text-xs font-medium text-white'
        >
          {label}
        </div>
      </EdgeLabelRenderer>
    </>
  );
};

const edgeTypes = {
  custom: CustomEdge,
  // You can add other custom edge components here
};

export default edgeTypes;
