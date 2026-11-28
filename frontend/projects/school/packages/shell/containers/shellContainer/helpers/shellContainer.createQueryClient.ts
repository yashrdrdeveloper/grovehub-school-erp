import { QueryClient } from '@tanstack/react-query';
import { SHELL_QUERY_OPTIONS } from '../constants/shellContainer.queryOptions';
import type { CreateShellQueryClient } from './types/shellContainer.createQueryClient';

export const createShellQueryClient: CreateShellQueryClient = () =>
  new QueryClient(SHELL_QUERY_OPTIONS);
