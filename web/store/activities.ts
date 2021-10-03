import { Work } from './works'

export const ActivityType = {
  NEW: 1,
  UPDATE: 2,
} as const
// eslint-disable-next-line no-redeclare
export type ActivityType = typeof ActivityType[keyof typeof ActivityType]

export interface User {
  id: string
  name: string
  nickname: string
  picture: string
}

export interface Activity {
  id: number
  type: ActivityType
  user: User
  work: Work
  createdAt: Date
}

export interface State {
  activities: Activity[]
}

export const state: () => State = () => ({
  activities: [
    {
      id: 0,
      type: 1,
      user: {
        id: 'aaaaa',
        name: 'XXX',
        nickname: 'XXX',
        picture:
          'https://cdn.pixabay.com/photo/2021/05/02/08/33/jellyfish-6222849_960_720.jpg',
      },
      work: {
        id: '000',
        author: 'abc',
        title: 'YYY',
        contentUrl: '',
        description: '',
        thumbnailUrl: '',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      createdAt: new Date(),
    },
    {
      id: 1,
      type: 2,
      user: {
        id: 'aaaaa',
        name: 'XXX',
        nickname: 'XXX',
        picture:
          'https://cdn.pixabay.com/photo/2021/05/02/08/33/jellyfish-6222849_960_720.jpg',
      },
      work: {
        id: '000',
        author: 'abc',
        title: 'YYY',
        contentUrl: '',
        description: '',
        thumbnailUrl: '',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      createdAt: new Date(),
    },
  ],
})
