import { readFileSync } from 'fs';
import { join } from 'path';

export interface Activity {
  name: string;
  color?: string;
  dates: string[];
}

export interface ActivitiesData {
  activities: Record<string, Activity>;
}

export interface ProcessedActivityData {
  name: string;
  color?: string;
  dateSet: Set<string>;
  totalCount: number;
}

export function loadActivitiesData(filePath: string = 'data/activities.json'): ActivitiesData {
  try {
    const fullPath = join(process.cwd(), filePath);
    const fileContent = readFileSync(fullPath, 'utf-8');
    return JSON.parse(fileContent) as ActivitiesData;
  } catch (error) {
    console.error('Error loading activities data:', error);
    // Return default empty data if file doesn't exist
    return { activities: {} };
  }
}

export function processActivityData(rawData: ActivitiesData): Record<string, ProcessedActivityData> {
  const processed: Record<string, ProcessedActivityData> = {};
  
  for (const [key, activity] of Object.entries(rawData.activities)) {
    processed[key] = {
      name: activity.name,
      color: activity.color,
      dateSet: new Set(activity.dates),
      totalCount: activity.dates.length,
    };
  }
  
  return processed;
}

export function getActivityKeys(data: ActivitiesData): string[] {
  return Object.keys(data.activities);
}

export function getActivityLevel(date: string, activityData: ProcessedActivityData): number {
  if (!activityData.dateSet.has(date)) {
    return 0;
  }
  
  // For now, just return 1 for any activity
  // In the future, you could track multiple occurrences per day
  return 1;
}

export function getActivityLevelForMultiple(date: string, activities: ProcessedActivityData[]): number {
  let count = 0;
  for (const activity of activities) {
    if (activity.dateSet.has(date)) {
      count++;
    }
  }
  
  // Convert to level (0-4)
  if (count === 0) return 0;
  if (count === 1) return 1;
  if (count === 2) return 2;
  if (count >= 3) return 3;
  return 4;
}