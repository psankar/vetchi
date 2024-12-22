import { Opening } from './opening';

export interface Application {
  id: string;
  opening: Opening;
  state: ApplicationState;
  resume_url: string;
  cover_letter?: string;
  created_at: string;
  last_updated_at: string;
}

export enum ApplicationState {
  SUBMITTED_APPLICATION = 'SUBMITTED_APPLICATION',
  SHORTLISTED_APPLICATION = 'SHORTLISTED_APPLICATION',
  REJECTED_APPLICATION = 'REJECTED_APPLICATION',
  WITHDRAWN_APPLICATION = 'WITHDRAWN_APPLICATION',
}

export interface CreateApplicationRequest {
  opening_id: string;
  resume_url: string;
  cover_letter?: string;
}

export interface FilterApplicationsRequest {
  state?: ApplicationState[];
  from_date?: string;
  to_date?: string;
  opening_id?: string;
  pagination_key?: string;
  limit?: number;
}

export interface GetApplicationRequest {
  id: string;
}

export interface WithdrawApplicationRequest {
  id: string;
}

export interface ApplicationMessage {
  id: string;
  application_id: string;
  sender_id: string;
  sender_type: 'HUB_USER' | 'ORG_USER';
  message: string;
  created_at: string;
}

export interface SendMessageRequest {
  application_id: string;
  message: string;
}

export interface GetMessagesRequest {
  application_id: string;
  pagination_key?: string;
  limit?: number;
} 