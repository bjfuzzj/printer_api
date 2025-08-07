Option Strict Off
Option Explicit On
Imports VB = Microsoft.VisualBasic

Friend Class Form1
	Inherits System.Windows.Forms.Form
	' =====================================================
	'
    '   .Net VB Status API Sample Program
	'
	'
	' =====================================================
	
	Dim TxCnt As Short
    Dim CY As Integer

	
    Private Sub Form1_Load(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles MyBase.Load

        Text2.Text = ""

        '********************* Initialize�iSelect Port Number�j
        CY = PortInitialize("USB001")
        If CY < 0 Then MsgBox("Error!", vbOKOnly) : End

    End Sub


    '********************* Get Version Infomation
    Private Sub Button1_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button1.Click
        Dim s As String = New String("", 255)
        Dim l As Integer

        l = GetFirmwVersion(CY, s)
        If l >= 0 Then Text2.Text = VB.Left(s, l) Else Text2.Text = "ERROR!"

    End Sub
	
    '********************* Printer Status
    Private Sub Button2_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button2.Click
        Dim stat As Integer

        stat = GetStatus(CY)

        If stat = STATUS_ERROR Then
            Text2.Text = "ERROR!"
            Exit Sub
        End If

        If stat And GROUP_USUALLY Then '================ Usually status group.

            Select Case stat
                Case STATUS_USUALLY_IDLE : Text2.Text = "Idle"
                Case STATUS_USUALLY_PRINTING : Text2.Text = "Printing"
                Case STATUS_USUALLY_STANDSTILL : Text2.Text = "STANDSTILL"
                Case STATUS_USUALLY_PAPER_END : Text2.Text = "Paper End"
                Case STATUS_USUALLY_RIBBON_END : Text2.Text = "Ribbon End"
                Case STATUS_USUALLY_COOLING : Text2.Text = "Head Cooling Down"
                Case STATUS_USUALLY_MOTCOOLING : Text2.Text = "Motor Cooling Down"
            End Select

        ElseIf stat And GROUP_SETTING Then  '============ Printer setting status group.

            Select Case stat
                Case STATUS_SETTING_COVER_OPEN : Text2.Text = "Cover Open"
                Case STATUS_SETTING_PAPER_JAM : Text2.Text = "Paper Jam"
                Case STATUS_SETTING_RIBBON_ERR : Text2.Text = "Ribbon Error"
                Case STATUS_SETTING_PAPER_ERR : Text2.Text = "Paper definition Error"
                Case STATUS_SETTING_DATA_ERR : Text2.Text = "Data Error"
                Case STATUS_SETTING_SCRAPBOX_ERR : Text2.Text = "Scrap Box Error"
            End Select

        ElseIf stat And GROUP_HARDWARE Then  '=========== Hardware erro status group.

            Select Case stat
                Case STATUS_HARDWARE_ERR01 : Text2.Text = "Head Voltage Error"
                Case STATUS_HARDWARE_ERR02 : Text2.Text = "Head Position Error"
                Case STATUS_HARDWARE_ERR03 : Text2.Text = "Fan Stop Error"
                Case STATUS_HARDWARE_ERR04 : Text2.Text = "Cutter Error"
                Case STATUS_HARDWARE_ERR05 : Text2.Text = "Pinch Roller Error"
                Case STATUS_HARDWARE_ERR06 : Text2.Text = "Illegal Head Temperature"
                Case STATUS_HARDWARE_ERR07 : Text2.Text = "Illegal Media Temperature"
                Case STATUS_HARDWARE_ERR08 : Text2.Text = "Ribbon Tension Error"
                Case STATUS_HARDWARE_ERR09 : Text2.Text = "RFID Module Error"
                Case STATUS_HARDWARE_ERR10 : Text2.Text = "Illegal Motor Temperature"
            End Select
        ElseIf stat And GROUP_SYSTEM Then  '============= System error status group.

            Text2.Text = "SYSTEM ERROR"

        ElseIf stat And GROUP_FLSHPROG Then  '=========== Flash rewriting staus group.

            Text2.Text = "FLSHPROG MODE"

        End If

    End Sub

    '********************* Get Serial Number
    Private Sub Button3_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button3.Click
        Dim s As String = New String("", 255)
        Dim l As Integer

        l = GetSerialNo(CY, s)
        If l >= 0 Then Text2.Text = VB.Left(s, l) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Get Sensor Information
    Private Sub Button4_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button4.Click
        Dim s As String = New String("", 255)
        Dim l As Integer

        l = GetSensorInfo(CY, s)
        If l >= 0 Then Text2.Text = VB.Left(s, l) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Get Color data Version
    Private Sub Button5_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button5.Click
        Dim s As String = New String("", 255)
        Dim ver As String
        Dim cs As String
        Dim l As Integer

        l = GetColorDataVersion(CY, s)
        If l >= 0 Then ver = VB.Left(s, l) Else Text2.Text = "ERROR!" : Exit Sub

        l = GetColorDataChecksum(CY, s)
        If l >= 0 Then cs = VB.Left(s, l) Else Text2.Text = "ERROR!" : Exit Sub

        Text2.Text = ver + "   " + cs

    End Sub

    '********************* Get Life counter
    Private Sub Button7_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button7.Click
        Dim c As Integer

        c = GetCounterL(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Get Counter A
    Private Sub Button8_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button8.Click
        Dim c As Integer

        c = GetCounterA(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Clear Counter A
    Private Sub Button9_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button9.Click

        If SetClearCounterA(CY) <> 0 Then Text2.Text = "True" Else Text2.Text = "False"

    End Sub

    '********************* Get Counter B
    Private Sub Button10_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button10.Click
        Dim c As Integer

        c = GetCounterB(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Clear Counter B
    Private Sub Button11_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button11.Click

        If SetClearCounterB(CY) <> 0 Then Text2.Text = "True" Else Text2.Text = "False"

    End Sub

    '********************** Get counter P
    Private Sub Button12_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button12.Click
        Dim c As Long

        c = GetCounterP(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************** Set counter P
    Private Sub Button13_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button13.Click

        If SetCounterP(CY, CInt(Text3.Text)) <> 0 Then Text2.Text = "True" Else Text2.Text = "False"

    End Sub

    '********************** Get counter Matte
    Private Sub Button14_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button14.Click
        Dim c As Integer

        c = GetCounterMatte(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************** Get counter M
    Private Sub Button15_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button15.Click
        Dim c As Integer

        c = GetCounterM(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************** Clear counter M
    Private Sub Button16_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button16.Click

        If SetClearCounterM(CY) <> 0 Then Text2.Text = "True" Else Text2.Text = "False"

    End Sub

    '********************* Get Remaining print quantity
    Private Sub Button17_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button17.Click
        Dim c As Integer

        c = GetPQTY(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Get the number of free image buffer
    Private Sub Button18_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button18.Click
        Dim c As Integer

        c = GetFreeBuffer(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Get Resolution H
    Private Sub Button19_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button19.Click
        Dim c As Integer

        c = GetResolutionH(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Get Resolution V
    Private Sub Button20_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button20.Click
        Dim c As Integer

        c = GetResolutionV(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Set Cutter non scarp mode
    Private Sub Button21_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button21.Click

        If SetCutterMode(CY, CUTTER_MODE_NONSCRAP) >= 0 Then Text2.Text = "True" Else Text2.Text = "False"

    End Sub

    '********************* Get media code
    Private Sub Button22_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button22.Click
        Dim s As String = New String("", 255)
        Dim l As Integer

        l = GetMedia(CY, s)
        If l >= 0 Then Text2.Text = VB.Left(s, l) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Get Media counter
    Private Sub Button23_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button23.Click
        Dim c As Integer

        c = GetMediaCounter(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Get Media lot information
    Private Sub Button24_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button24.Click
        Dim s As String = New String("", 255)
        Dim l As Integer

        l = GetMediaLotNo(CY, s)
        If l >= 0 Then Text2.Text = VB.Left(s, l) Else Text2.Text = "ERROR!"

    End Sub

    '******************** Get Media color offset value
    Private Sub Button25_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button25.Click
        Dim s As String = New String("", 255)
        Dim l As Integer

        l = GetMediaColorOffset(CY)
        If l >= 0 Then
            s = VB.Right("0" & Hex(l), 8)
            Text2.Text = s
        Else
            Text2.Text = "ERROR!"
        End If

    End Sub

    '********************* Get Media ID setting value
    Private Sub Button26_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button26.Click
        Dim c As Long

        c = GetMediaIdSetInfo(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"

    End Sub

    '********************* Get Media mounter default value
    Private Sub Button27_Click(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles Button27.Click
        Dim c As Long

        c = GetInitialMediaCount(CY)
        If c >= 0 Then Text2.Text = Str(c) Else Text2.Text = "ERROR!"
    End Sub

End Class