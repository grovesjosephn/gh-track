import React, { useState } from 'react';
import { Box, Text } from 'ink';
import { ContributionsGrid } from './components/ContributionsGrid.js';
import { generateContributionsGrid } from './utils/dateUtils.js';
import { loadActivitiesData, processActivityData } from './utils/dataLoader.js';

function App() {
  const [grid] = useState(() => generateContributionsGrid());
  const [activitiesData] = useState(() => {
    const rawData = loadActivitiesData();
    return processActivityData(rawData);
  });
  
  const activityKeys = Object.keys(activitiesData);

  if (activityKeys.length === 0) {
    return (
      <Box flexDirection="column" padding={1}>
        <Text color="red">ERROR: No activities found. Please check your data/activities.json file.</Text>
      </Box>
    );
  }

  return (
    <Box flexDirection="column" padding={1}>
      <Box marginBottom={2}>
        <Text color="cyan" bold>
          Activity Tracker - All Activities
        </Text>
      </Box>
      
      {activityKeys.map((activityKey, index) => {
        const activity = activitiesData[activityKey];
        return (
          <Box key={activityKey} flexDirection="column" marginBottom={index < activityKeys.length - 1 ? 2 : 0}>
            <ContributionsGrid
              grid={grid}
              activityData={activity}
              title={activity.name}
              activityColor={activity.color}
            />
          </Box>
        );
      })}
    </Box>
  );
}

export default App;
