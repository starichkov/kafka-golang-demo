# 🏷️ Go Test & Coverage Commands Cheat Sheet

## 🚦 Run All Tests (Verbose Output)

```bash
go test -v ./...
```

## 🟢 Run Only Unit Tests in One Package

```bash
go test -v ./internal/kafka
```

## 🧪 Run Tests With the Race Detector

```bash
go test -race ./...
```

## 🟩 Run All Tests and Generate Coverage Report

```bash
go test -coverprofile=coverage.out -covermode=atomic ./...
```

## 🖼️ View Coverage Report in Browser

```bash
go tool cover -html=coverage.out
```

## 📊 See Coverage for Each Function

```bash
go tool cover -func=coverage.out
```

## 🧹 Combine Multiple Coverage Profiles

```bash
echo "mode: set" > total_coverage.out
tail -n +2 unit_coverage.out >> total_coverage.out
tail -n +2 integration_coverage.out >> total_coverage.out
# Then view or upload `total_coverage.out`
```

## 🧭 List All Go Packages in the Module

```bash
go list ./...
```

## 🔎 See Which Tests Are Detected and Run

```bash
go test -list . ./...
```

## ⚡ Run Only Integration Tests (with tags, if needed)

```bash
go test -tags=integration ./...
```

## 🛠️ (Optional) Test Only One Function or File

```bash
go test -run TestMyFunc ./internal/kafka
go test ./internal/kafka/consumer_test.go
```

---

**Tip:**  
Add these to your project’s README or `TESTING.md` for easy reference!
