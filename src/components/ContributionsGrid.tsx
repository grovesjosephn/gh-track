import React from "react";
import { Box, Text } from "ink";
import { formatDate } from "../utils/dateUtils.js";
import type { ProcessedActivityData } from "../utils/dataLoader.js";
import { getActivityLevel } from "../utils/dataLoader.js";

interface ContributionsGridProps {
  grid: Date[][];
  activityData?: ProcessedActivityData;
  selectedWeek?: number;
  selectedDay?: number;
  activityColor?: string;
  title?: string;
}

const getActivityColor = (color: string | undefined, level: number): string => {
  if (level === 0) return "gray";
  return color || "green";
};

const LEVEL_CHARS = {
  0: "▒", // Light shade
  1: "░", // Middle shade
  2: "▓", // Dark shade
  3: "█", // Full block
  4: "█", // Full block
} as const;


function getCharForLevel(level: number): string {
  return LEVEL_CHARS[level as keyof typeof LEVEL_CHARS] || "▒";
}

export function ContributionsGrid({
  grid,
  activityData,
  selectedWeek = -1,
  selectedDay = -1,
  activityColor,
  title,
}: ContributionsGridProps) {
  const dayLabels = ["S", "M", "T", "W", "T", "F", "S"];

  return (
    <Box flexDirection="column" paddingX={1}>
      {/* Activity title */}
      {title && (
        <Box marginBottom={1}>
          <Text bold color={activityColor || activityData?.color || "cyan"}>
            {title} ({activityData?.totalCount || 0} activities)
          </Text>
        </Box>
      )}

      {/* Grid rows with day labels */}
      {dayLabels.map((dayLabel, dayIndex) => (
        <Box key={dayIndex}>
          <Box width={2} justifyContent="center">
            <Text color="dim" bold>
              {dayLabel}
            </Text>
          </Box>

          {grid.map((week, weekIndex) => {
            const date = week[dayIndex];
            const isValidDate = date && date.getTime() > 0;
            const dateString = isValidDate ? formatDate(date) : "";
            const level =
              isValidDate && activityData
                ? getActivityLevel(dateString, activityData)
                : 0;
            const isSelected =
              weekIndex === selectedWeek && dayIndex === selectedDay;

            const effectiveColor = activityColor || activityData?.color;

            return (
              <Box key={`${weekIndex}-${dayIndex}`} width={2} justifyContent="center">
                <Text
                  color={isSelected ? "inverse" : getActivityColor(effectiveColor, level)}
                  backgroundColor={isSelected ? "cyan" : undefined}
                >
                  {isValidDate ? getCharForLevel(level) : " "}
                </Text>
              </Box>
            );
          })}
        </Box>
      ))}

      {/* Legend */}
      <Box marginTop={1} alignItems="center">
        <Box width={2}></Box>
        <Text color="dim">Less </Text>
        {[0, 1, 2, 3, 4].map((level) => (
          <Text key={level}>{getCharForLevel(level)}</Text>
        ))}
        <Text color="dim"> More</Text>
      </Box>

    </Box>
  );
}
