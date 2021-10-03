export interface RootState {
  version: string
}

export const state: () => RootState = () => ({
  version: '1.0.0',
})
