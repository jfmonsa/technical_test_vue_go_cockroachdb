/**
 * Represents a stock recommendation or analysis entry.
 *
 */
export interface Stock {
  ticker: string;
  company: string;
  brokerage: string;
  action: string;
  rating_from: string;
  rating_to: string;
  target_from: number;
  target_to: number;
  time: string;
}
