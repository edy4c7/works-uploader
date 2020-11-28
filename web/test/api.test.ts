// eslint-disable-next-line import/no-duplicates
import axios from 'axios'
// eslint-disable-next-line import/no-duplicates
import { AxiosResponse } from 'axios'
import { ApiClientImpl, Work } from '~/plugins/api'

jest.mock('axios')

describe('ApiImpl', () => {
  it('getWorks', async () => {
    const data: Work[] = [
      {
        id: '001',
        author: 'taro',
        title: 'work01',
        description: 'hoge',
        thumbnailUrl: 'https://example.com/work01',
        contentUrl: 'https://example.com/work01',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      {
        id: '002',
        author: 'jiro',
        title: 'work02',
        description: 'fuga',
        thumbnailUrl: 'https://example.com/work02',
        contentUrl: 'https://example.com/work02',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
    ]
    const res: AxiosResponse<Work[]> = {
      data,
      status: 200,
      statusText: '',
      headers: '',
      config: {},
    }

    ;(axios.get as jest.Mock).mockResolvedValue(res)
    const api = new ApiClientImpl(axios)
    expect(await api.getWorks()).toEqual(data)
  })
})
