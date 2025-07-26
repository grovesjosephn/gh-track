export interface DateRange {
  start: Date;
  end: Date;
}

export function getYearDateRange(year?: number): DateRange {
  const currentYear = year || new Date().getFullYear();
  const start = new Date(currentYear, 0, 1); // January 1st
  const end = new Date(currentYear, 11, 31); // December 31st
  
  return { start, end };
}

export function getContributionsDateRange(): DateRange {
  const today = new Date();
  const end = new Date(today);
  const start = new Date(today);
  start.setDate(start.getDate() - 364); // 365 days total (including today)
  
  return { start, end };
}

export function formatDate(date: Date): string {
  return date.toISOString().split('T')[0];
}

export function parseDate(dateString: string): Date {
  return new Date(dateString + 'T00:00:00.000Z');
}

export function getWeekNumber(date: Date): number {
  const start = new Date(date.getFullYear(), 0, 1);
  const days = Math.floor((date.getTime() - start.getTime()) / (24 * 60 * 60 * 1000));
  return Math.floor(days / 7);
}

export function getDayOfWeek(date: Date): number {
  return date.getDay(); // 0 = Sunday, 1 = Monday, etc.
}

export function getAllDatesInRange(start: Date, end: Date): Date[] {
  const dates: Date[] = [];
  const current = new Date(start);
  
  while (current <= end) {
    dates.push(new Date(current));
    current.setDate(current.getDate() + 1);
  }
  
  return dates;
}

export function getWeekStartDate(date: Date): Date {
  const start = new Date(date);
  const day = start.getDay();
  const diff = start.getDate() - day; // Get Sunday
  return new Date(start.setDate(diff));
}

export function generateContributionsGrid(): Date[][] {
  const { start, end } = getContributionsDateRange();
  const allDates = getAllDatesInRange(start, end);
  
  // Create a 7x53 grid (7 days x ~53 weeks)
  const grid: Date[][] = [];
  for (let week = 0; week < 53; week++) {
    grid[week] = [];
    for (let day = 0; day < 7; day++) {
      grid[week][day] = new Date(0); // placeholder
    }
  }
  
  // Fill the grid with actual dates
  let weekIndex = 0;
  let currentWeekStart = getWeekStartDate(start);
  
  for (const date of allDates) {
    const weekStart = getWeekStartDate(date);
    
    // Move to next week if needed
    if (weekStart.getTime() !== currentWeekStart.getTime()) {
      weekIndex++;
      currentWeekStart = weekStart;
    }
    
    const dayOfWeek = getDayOfWeek(date);
    if (weekIndex < 53) {
      grid[weekIndex][dayOfWeek] = date;
    }
  }
  
  return grid;
}