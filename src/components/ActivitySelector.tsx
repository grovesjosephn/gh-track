import React from 'react';
import { Box, Text } from 'ink';
import type { ProcessedActivityData } from '../utils/dataLoader.js';

interface ActivitySelectorProps {
  activities: Record<string, ProcessedActivityData>;
  selectedActivity: string;
  onActivityChange: (activityKey: string) => void;
}

export function ActivitySelector({ activities, selectedActivity }: ActivitySelectorProps) {
  const activityKeys = Object.keys(activities);
  const currentIndex = activityKeys.indexOf(selectedActivity);
  const currentActivity = activities[selectedActivity];
  
  if (activityKeys.length === 0) {
    return (
      <Box flexDirection="column" marginBottom={2} paddingX={1} borderStyle="round" borderColor="red">
        <Text color="red">ERROR: No activities found. Please check your data/activities.json file.</Text>
      </Box>
    );
  }
  
  return (
    <Box flexDirection="column" marginBottom={2} paddingX={1}>
      {/* Title and current activity */}
      <Box marginBottom={1} justifyContent="space-between">
        <Text color="cyan" bold>
          Activity Tracker - {currentActivity?.name || 'Unknown Activity'}
        </Text>
        <Text color="dim">
          ({currentIndex + 1}/{activityKeys.length})
        </Text>
      </Box>
      
      {/* Activity stats */}
      <Box marginBottom={1}>
        <Text color="green">
          Total activities: {currentActivity ? currentActivity.totalCount : 0}
        </Text>
      </Box>
      
      {/* Navigation tabs */}
      <Box marginBottom={1}>
        {activityKeys.map((key) => {
          const activity = activities[key];
          const isSelected = key === selectedActivity;
          return (
            <Box key={key} marginRight={1}>
              <Text
                color={isSelected ? 'cyan' : 'dim'}
                backgroundColor={isSelected ? 'cyan' : undefined}
                inverse={isSelected}
              >
                {isSelected ? `[${activity.name}]` : ` ${activity.name} `}
              </Text>
            </Box>
          );
        })}
      </Box>
      
      {/* Controls hint */}
      <Box borderStyle="round" borderColor="blue" paddingX={1}>
        <Box flexDirection="column">
          <Text color="blue" bold>Controls:</Text>
          <Text color="dim">← → Switch activities (when not in grid)  •  ↑ ↓ W S Enter/navigate grid  •  A D Navigate horizontally  •  SPACE Exit grid  •  Q Quit</Text>
        </Box>
      </Box>
    </Box>
  );
}