import { Box, Card, CardContent, Grid, Typography } from '@mui/material'

export default function Dashboard() {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        ダッシュボード
      </Typography>
      <Grid container spacing={3}>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                総プロジェクト数
              </Typography>
              <Typography variant="h4">-</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card sx={{ bgcolor: 'error.light', color: 'white' }}>
            <CardContent>
              <Typography gutterBottom>遅延プロジェクト</Typography>
              <Typography variant="h4">-</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card sx={{ bgcolor: 'warning.light', color: 'white' }}>
            <CardContent>
              <Typography gutterBottom>注意プロジェクト</Typography>
              <Typography variant="h4">-</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card sx={{ bgcolor: 'success.light', color: 'white' }}>
            <CardContent>
              <Typography gutterBottom>正常プロジェクト</Typography>
              <Typography variant="h4">-</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                組織別プロジェクト状況
              </Typography>
              <Typography variant="body2" color="textSecondary">
                データがありません。バックエンドAPIと連携後に表示されます。
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  )
}
