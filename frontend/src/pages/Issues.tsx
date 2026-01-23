import { Box, Card, CardContent, Typography } from '@mui/material'

export default function Issues() {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        チケット一覧
      </Typography>
      <Card>
        <CardContent>
          <Typography variant="body1">
            遅延チケットの一覧とフィルタ機能を提供します。
          </Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mt: 2 }}>
            この機能はFRONT-005チケットで実装予定です。
          </Typography>
        </CardContent>
      </Card>
    </Box>
  )
}
