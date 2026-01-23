import { Box, Card, CardContent, Typography } from '@mui/material'

export default function Projects() {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        プロジェクト一覧
      </Typography>
      <Card>
        <CardContent>
          <Typography variant="body1">
            プロジェクトの一覧と遅延状況を表示します。
          </Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mt: 2 }}>
            この機能はFRONT-004チケットで実装予定です。
          </Typography>
        </CardContent>
      </Card>
    </Box>
  )
}
