import { Box, Card, CardContent, Typography } from '@mui/material'

export default function Organizations() {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        組織管理
      </Typography>
      <Card>
        <CardContent>
          <Typography variant="body1">
            組織階層の管理画面です。プロジェクトの組織紐付けを行うことができます。
          </Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mt: 2 }}>
            この機能はFRONT-007チケットで実装予定です。
          </Typography>
        </CardContent>
      </Card>
    </Box>
  )
}
