package msisdn

import (
	"bytes"
	"encoding/json"
	"github.com/andybalholm/brotli"
	"io/ioutil"
	"strings"
	"sync"
)

// Zone contains a single Country's Zone information
type Zone struct {
	Name string
}

// Country contains a single Country's information
type Country struct {
	Code     string
	Name     string
	Language string
	Zones    []Zone
}

// GENERATED FILE DO NOT MODIFY DIRECTLY
var (
	onceCountry sync.Once
	mapped      map[string]Country
	data        = []byte{139, 135, 47, 17, 85, 178, 6, 64, 235, 3, 155, 200, 28, 79, 234, 235, 172, 44, 8, 157, 61, 126, 161, 243, 83, 176, 185, 180, 89, 106, 126, 16, 74, 251, 192, 157, 27, 163, 77, 80, 95, 90, 246, 90, 46, 99, 254, 113, 74, 52, 13, 29, 136, 4, 80, 66, 23, 129, 65, 113, 27, 253, 52, 1, 174, 111, 238, 75, 213, 182, 53, 203, 10, 168, 154, 98, 189, 192, 242, 217, 91, 218, 156, 111, 22, 157, 174, 62, 102, 106, 146, 85, 171, 73, 119, 111, 202, 82, 3, 62, 166, 104, 2, 43, 121, 1, 210, 145, 250, 54, 143, 211, 227, 186, 219, 36, 70, 181, 37, 76, 93, 120, 154, 74, 107, 104, 94, 185, 15, 192, 183, 147, 38, 245, 250, 217, 151, 196, 98, 223, 155, 170, 173, 157, 75, 55, 77, 200, 165, 74, 207, 236, 142, 81, 252, 25, 93, 229, 152, 74, 83, 149, 157, 113, 255, 1, 10, 41, 65, 249, 232, 196, 15, 199, 211, 64, 119, 228, 87, 130, 34, 79, 41, 130, 115, 51, 6, 200, 83, 56, 242, 20, 0, 210, 145, 154, 222, 101, 169, 20, 124, 25, 243, 253, 255, 90, 197, 113, 59, 254, 194, 156, 32, 202, 10, 1, 219, 190, 255, 251, 253, 241, 76, 15, 191, 63, 40, 69, 15, 31, 207, 244, 224, 97, 216, 235, 18, 128, 53, 189, 243, 224, 193, 234, 27, 139, 177, 239, 255, 183, 254, 29, 232, 122, 248, 175, 239, 135, 235, 179, 158, 248, 78, 117, 166, 155, 199, 19, 177, 1, 181, 14, 199, 111, 128, 204, 111, 217, 124, 249, 2, 24, 141, 244, 105, 195, 128, 155, 247, 20, 145, 149, 232, 76, 40, 67, 174, 244, 176, 128, 209, 217, 50, 158, 96, 82, 41, 17, 210, 105, 49, 220, 74, 126, 152, 196, 179, 65, 56, 87, 97, 31, 41, 17, 185, 13, 39, 252, 11, 33, 140, 168, 99, 175, 167, 51, 48, 120, 212, 253, 243, 121, 35, 76, 160, 137, 120, 83, 134, 126, 121, 62, 21, 141, 31, 63, 49, 47, 81, 11, 44, 164, 164, 7, 219, 162, 192, 35, 87, 52, 238, 229, 15, 237, 82, 130, 117, 34, 3, 43, 106, 120, 6, 249, 112, 238, 163, 64, 56, 254, 223, 142, 20, 199, 77, 202, 161, 135, 251, 208, 209, 203, 136, 32, 69, 234, 202, 218, 151, 112, 77, 50, 125, 216, 234, 93, 183, 153, 39, 236, 218, 178, 93, 216, 6, 68, 133, 53, 207, 88, 154, 62, 247, 31, 163, 110, 79, 90, 206, 179, 79, 142, 14, 222, 6, 16, 211, 86, 130, 183, 247, 32, 213, 189, 79, 81, 185, 12, 174, 236, 240, 86, 210, 111, 5, 114, 112, 73, 160, 3, 76, 84, 140, 161, 225, 5, 203, 21, 39, 52, 97, 22, 35, 194, 174, 135, 223, 31, 96, 192, 101, 108, 126, 37, 127, 43, 161, 235, 249, 5, 0, 182, 228, 95, 40, 82, 122, 57, 98, 212, 97, 183, 207, 160, 21, 25, 194, 212, 172, 199, 118, 120, 84, 100, 222, 15, 9, 190, 12, 211, 139, 60, 243, 193, 218, 15, 45, 9, 14, 251, 150, 136, 177, 202, 192, 210, 130, 112, 157, 90, 190, 26, 91, 162, 210, 57, 147, 231, 171, 23, 67, 28, 237, 203, 240, 123, 61, 117, 95, 66, 218, 184, 15, 198, 74, 151, 47, 242, 212, 35, 7, 77, 164, 129, 233, 135, 223, 248, 115, 128, 181, 93, 104, 80, 252, 155, 235, 242, 43, 116, 61, 176, 115, 163, 192, 184, 152, 71, 156, 233, 111, 176, 57, 3, 177, 28, 25, 53, 168, 16, 20, 8, 235, 172, 248, 233, 52, 233, 162, 155, 248, 165, 150, 21, 115, 18, 12, 60, 133, 251, 220, 227, 173, 21, 94, 116, 225, 149, 53, 56, 83, 226, 41, 94, 129, 78, 224, 41, 81, 186, 8, 108, 7, 69, 222, 22, 130, 65, 217, 128, 82, 99, 0, 77, 8, 40, 185, 218, 46, 223, 43, 226, 29, 35, 96, 47, 191, 164, 23, 101, 204, 151, 163, 184, 27, 108, 5, 207, 16, 149, 90, 218, 119, 223, 94, 222, 40, 153, 45, 215, 149, 8, 168, 160, 70, 31, 48, 70, 134, 116, 37, 75, 121, 89, 86, 42, 84, 146, 57, 58, 29, 121, 45, 161, 17, 100, 170, 88, 51, 85, 44, 172, 246, 12, 63, 214, 187, 108, 46, 98, 232, 193, 39, 46, 193, 192, 206, 171, 238, 175, 204, 31, 203, 253, 76, 29, 186, 91, 23, 121, 173, 35, 28, 237, 71, 135, 236, 102, 36, 55, 246, 60, 31, 19, 245, 88, 77, 194, 140, 4, 228, 182, 89, 189, 160, 190, 37, 88, 126, 161, 220, 144, 85, 154, 252, 168, 166, 22, 22, 164, 168, 151, 142, 122, 89, 126, 111, 197, 81, 180, 188, 194, 173, 236, 219, 107, 191, 156, 125, 50, 132, 125, 240, 240, 36, 3, 211, 165, 176, 0, 104, 199, 139, 55, 13, 164, 149, 116, 228, 64, 122, 128, 28, 182, 124, 103, 169, 44, 125, 100, 95, 246, 65, 233, 54, 33, 26, 191, 119, 198, 29, 126, 76, 237, 234, 146, 168, 211, 231, 206, 130, 174, 49, 223, 2, 167, 54, 170, 103, 235, 7, 37, 171, 159, 182, 239, 8, 166, 167, 177, 48, 104, 204, 114, 19, 55, 222, 151, 180, 225, 182, 97, 42, 251, 37, 157, 29, 108, 140, 29, 44, 180, 45, 112, 48, 160, 204, 190, 236, 220, 7, 83, 131, 13, 13, 24, 153, 179, 56, 48, 82, 0, 195, 240, 139, 105, 200, 209, 119, 223, 170, 152, 230, 12, 71, 89, 52, 219, 103, 3, 166, 141, 123, 160, 19, 215, 11, 224, 98, 137, 63, 124, 149, 20, 6, 50, 116, 54, 63, 233, 202, 219, 192, 114, 167, 16, 61, 53, 174, 125, 165, 192, 72, 76, 234, 71, 254, 171, 169, 244, 160, 34, 80, 16, 196, 105, 3, 138, 168, 233, 195, 35, 3, 132, 191, 43, 124, 128, 229, 31, 78, 58, 114, 76, 81, 112, 81, 123, 71, 238, 185, 7, 93, 62, 232, 123, 156, 248, 94, 182, 99, 166, 161, 251, 58, 130, 239, 168, 42, 207, 96, 4, 61, 153, 76, 10, 95, 171, 30, 100, 7, 203, 1, 175, 46, 167, 167, 109, 168, 95, 252, 103, 82, 119, 212, 242, 152, 82, 127, 208, 234, 17, 97, 164, 14, 150, 17, 204, 109, 246, 229, 22, 92, 97, 176, 157, 230, 40, 108, 127, 173, 165, 227, 158, 90, 146, 117, 226, 144, 43, 140, 221, 152, 58, 14, 197, 220, 11, 53, 3, 219, 148, 30, 173, 208, 242, 153, 125, 76, 91, 142, 243, 242, 94, 126, 141, 73, 170, 22, 122, 80, 127, 69, 43, 172, 162, 207, 183, 53, 194, 96, 164, 95, 22, 88, 50, 207, 214, 109, 232, 230, 74, 61, 195, 143, 72, 107, 156, 190, 212, 6, 221, 47, 247, 138, 253, 244, 43, 184, 81, 194, 116, 190, 96, 203, 239, 237, 104, 234, 173, 175, 82, 34, 153, 183, 182, 149, 157, 49, 217, 184, 26, 126, 13, 130, 68, 40, 190, 84, 54, 227, 205, 238, 162, 12, 89, 106, 37, 53, 237, 36, 184, 200, 90, 110, 45, 127, 108, 129, 23, 209, 150, 232, 252, 173, 1, 175, 223, 121, 190, 88, 214, 195, 103, 169, 64, 169, 108, 220, 216, 101, 30, 18, 235, 26, 119, 209, 185, 193, 83, 133, 122, 125, 179, 127, 236, 218, 62, 191, 189, 210, 121, 117, 76, 72, 46, 31, 128, 5, 117, 236, 177, 53, 171, 242, 96, 87, 213, 40, 12, 118, 160, 128, 167, 227, 111, 177, 180, 136, 32, 124, 194, 251, 147, 117, 197, 85, 196, 200, 217, 194, 196, 156, 72, 23, 141, 148, 41, 10, 54, 9, 202, 150, 91, 145, 51, 218, 87, 210, 112, 159, 55, 176, 53, 250, 169, 84, 88, 26, 7, 203, 190, 249, 147, 20, 180, 208, 182, 218, 154, 231, 182, 99, 90, 226, 69, 134, 54, 228, 226, 131, 109, 58, 195, 187, 148, 15, 236, 194, 130, 38, 206, 189, 87, 9, 188, 189, 177, 0, 48, 45, 92, 131, 205, 165, 109, 137, 171, 50, 229, 133, 86, 220, 99, 223, 73, 31, 220, 133, 7, 109, 220, 89, 43, 242, 191, 168, 172, 37, 79, 183, 32, 115, 162, 211, 193, 179, 251, 251, 51, 31, 21, 141, 244, 251, 102, 89, 193, 168, 1, 12, 158, 77, 115, 142, 91, 102, 10, 136, 128, 170, 148, 160, 234, 140, 92, 135, 55, 58, 173, 188, 74, 197, 168, 96, 148, 202, 34, 91, 82, 202, 18, 138, 135, 192, 46, 200, 210, 44, 55, 54, 129, 43, 61, 240, 92, 132, 243, 225, 148, 25, 158, 35, 3, 8, 205, 66, 38, 232, 212, 160, 73, 101, 95, 208, 249, 228, 10, 167, 163, 92, 66, 195, 225, 196, 192, 197, 142, 213, 172, 229, 72, 11, 116, 2, 79, 10, 26, 82, 97, 102, 39, 188, 73, 228, 58, 45, 59, 206, 140, 45, 220, 40, 54, 139, 186, 153, 151, 221, 196, 234, 65, 138, 241, 60, 221, 92, 121, 169, 228, 86, 200, 45, 15, 221, 178, 182, 52, 91, 75, 110, 62, 213, 53, 81, 88, 93, 81, 174, 207, 2, 111, 66, 196, 79, 228, 35, 176, 84, 237, 117, 77, 6, 122, 24, 39, 54, 17, 82, 92, 237, 254, 28, 62, 229, 17, 157, 49, 84, 160, 20, 42, 95, 19, 112, 14, 24, 163, 106, 6, 190, 172, 19, 72, 159, 216, 200, 192, 50, 231, 210, 209, 117, 125, 61, 127, 66, 38, 26, 48, 220, 190, 58, 108, 185, 120, 30, 105, 70, 118, 161, 145, 119, 45, 192, 59, 234, 208, 15, 7, 19, 183, 144, 131, 159, 181, 27, 205, 231, 117, 255, 230, 204, 254, 178, 218, 39, 72, 164, 19, 85, 105, 85, 104, 61, 148, 103, 228, 211, 114, 173, 190, 38, 189, 140, 221, 224, 244, 208, 19, 245, 237, 88, 35, 248, 230, 119, 105, 52, 44, 49, 240, 210, 212, 185, 67, 228, 69, 103, 21, 79, 54, 159, 165, 242, 237, 159, 89, 227, 216, 144, 90, 205, 17, 6, 83, 114, 188, 83, 109, 203, 70, 142, 243, 83, 238, 36, 5, 237, 133, 7, 83, 107, 35, 162, 99, 72, 225, 93, 106, 225, 64, 6, 166, 181, 11, 108, 101, 170, 221, 46, 94, 244, 164, 50, 80, 24, 228, 101, 40, 214, 228, 181, 110, 231, 61, 186, 28, 160, 193, 233, 130, 43, 64, 19, 34, 224, 203, 41, 17, 173, 72, 33, 28, 40, 65, 58, 170, 126, 151, 156, 68, 200, 4, 156, 16, 80, 116, 108, 248, 159, 136, 10, 132, 91, 198, 108, 181, 251, 147, 50, 100, 81, 252, 201, 88, 173, 46, 64, 112, 93, 70, 220, 194, 233, 223, 84, 160, 156, 212, 111, 32, 150, 200, 185, 149, 242, 210, 193, 37, 70, 142, 13, 155, 45, 141, 51, 16, 124, 148, 126, 144, 160, 36, 85, 29, 146, 123, 206, 164, 128, 209, 239, 149, 222, 205, 94, 232, 33, 38, 69, 154, 121, 201, 131, 43, 40, 229, 139, 62, 248, 174, 183, 238, 167, 36, 66, 139, 190, 157, 132, 161, 40, 41, 19, 25, 218, 65, 187, 106, 240, 118, 74, 122, 59, 68, 19, 102, 114, 244, 236, 149, 132, 12, 45, 199, 20, 9, 133, 1, 109, 218, 121, 207, 24, 232, 63, 104, 191, 152, 165, 25, 241, 19, 77, 36, 3, 29, 60, 241, 213, 13, 66, 184, 166, 210, 187, 202, 155, 194, 91, 81, 115, 204, 78, 217, 55, 63, 62, 48, 9, 58, 9, 75, 98, 107, 17, 247, 237, 225, 66, 173, 178, 27, 250, 9, 7, 50, 48, 209, 221, 13, 152, 170, 231, 242, 20, 143, 214, 149, 160, 201, 193, 147, 9, 232, 225, 64, 35, 117, 29, 203, 255, 241, 156, 127, 30, 9, 40, 28, 104, 194, 20, 253, 203, 32, 135, 211, 12, 46, 75, 61, 175, 112, 201, 135, 42, 146, 34, 210, 252, 16, 107, 102, 32, 224, 21, 77, 164, 43, 89, 103, 99, 192, 44, 154, 216, 53, 184, 185, 138, 243, 167, 86, 6, 126, 65, 255, 8, 134, 5, 12, 119, 18, 2, 138, 170, 162, 173, 164, 207, 72, 149, 14, 3, 154, 58, 246, 100, 188, 83, 181, 50, 184, 168, 113, 41, 3, 9, 74, 60, 88, 45, 245, 248, 38, 71, 119, 130, 79, 250, 240, 59, 90, 228, 83, 107, 237, 181, 182, 82, 118, 50, 182, 9, 248, 14, 169, 42, 5, 181, 1, 188, 156, 81, 48, 66, 71, 7, 19, 13, 28, 91, 203, 215, 47, 186, 24, 121, 56, 255, 212, 189, 127, 176, 60, 240, 201, 209, 79, 22, 232, 109, 58, 93, 118, 203, 247, 14, 168, 87, 85, 162, 34, 81, 63, 7, 212, 86, 214, 80, 19, 9, 50, 49, 242, 209, 196, 40, 250, 11, 153, 128, 147, 162, 198, 20, 174, 232, 51, 120, 96, 79, 0, 75, 47, 190, 187, 229, 154, 156, 73, 120, 60, 50, 238, 168, 10, 21, 133, 177, 197, 194, 157, 255, 49, 141, 206, 36, 35, 19, 33, 53, 184, 251, 118, 127, 169, 183, 58, 83, 199, 68, 123, 219, 53, 190, 168, 64, 169, 111, 99, 115, 130, 177, 51, 197, 195, 140, 118, 198, 166, 33, 180, 23, 85, 168, 217, 139, 77, 39, 162, 64, 215, 124, 198, 42, 206, 197, 187, 152, 213, 172, 233, 176, 17, 62, 194, 8, 200, 23, 119, 17, 3, 15, 168, 29, 11, 163, 122, 44, 145, 219, 116, 188, 83, 9, 52, 201, 22, 69, 193, 100, 45, 59, 198, 165, 153, 21, 143, 105, 185, 238, 92, 234, 209, 92, 4, 225, 69, 47, 188, 238, 175, 249, 165, 178, 202, 41, 188, 125, 67, 208, 229, 204, 81, 125, 119, 176, 221, 96, 155, 200, 182, 104, 172, 237, 179, 241, 92, 231, 84, 22, 180, 113, 39, 83, 95, 122, 4, 109, 77, 34, 184, 20, 24, 141, 120, 24, 45, 204, 228, 214, 254, 0, 250, 36, 110, 62, 177, 170, 249, 25, 204, 233, 16, 73, 26, 14, 21, 119, 167, 56, 249, 78, 212, 52, 241, 23, 31, 128, 219, 178, 236, 254, 186, 143, 176, 204, 169, 106, 80, 52, 44, 74, 56, 239, 129, 208, 106, 200, 128, 179, 119, 166, 213, 191, 92, 49, 44, 16, 83, 186, 240, 74, 74, 9, 68, 161, 204, 83, 103, 5, 48, 221, 91, 188, 195, 79, 146, 82, 217, 18, 106, 52, 157, 230, 64, 21, 234, 245, 45, 113, 205, 109, 8, 215, 80, 135, 90, 113, 137, 82, 207, 59, 53, 104, 131, 90, 119, 111, 2, 246, 93, 60, 69, 87, 104, 77, 200, 139, 12, 44, 185, 196, 73, 44, 67, 255, 68, 131, 158, 229, 33, 52, 112, 100, 76, 158, 104, 1, 247, 112, 177, 12, 57, 122, 13, 71, 171, 6, 12, 78, 181, 155, 203, 122, 122, 250, 23, 77, 152, 21, 77, 128, 108, 229, 83, 42, 8, 133, 3, 45, 92, 7, 39, 254, 239, 240, 92, 87, 105, 186, 177, 210, 193, 115, 212, 156, 110, 161, 114, 199, 232, 190, 162, 158, 85, 3, 138, 129, 236, 165, 97, 242, 18, 17, 187, 207, 73, 7, 88, 128, 167, 155, 238, 184, 64, 142, 23, 107, 34, 65, 166, 68, 169, 199, 186, 102, 165, 27, 180, 5, 242, 241, 81, 179, 27, 204, 133, 8, 105, 176, 41, 52, 124, 129, 231, 32, 137, 250, 176, 128, 1, 0, 244, 10, 15, 42, 80, 10, 125, 235, 136, 104, 214, 90, 110, 230, 179, 242, 139, 55, 170, 142, 118, 229, 93, 116, 1, 53, 106, 190, 172, 18, 88, 126, 255, 161, 23, 180, 240, 170, 22, 63, 88, 190, 235, 212, 52, 195, 247, 92, 31, 50, 212, 169, 239, 70, 141, 27, 29, 178, 4, 175, 100, 140, 12, 196, 200, 39, 189, 193, 243, 98, 85, 107, 26, 35, 121, 96, 11, 128, 73, 80, 198, 102, 216, 157, 229, 149, 155, 92, 102, 106, 233, 9, 189, 64, 202, 91, 163, 101, 173, 106, 10, 195, 40, 130, 25, 9, 232, 28, 82, 212, 120, 200, 156, 242, 161, 60, 197, 132, 137, 74, 31, 125, 217, 227, 48, 252, 41, 225, 100, 5, 165, 176, 32, 35, 203, 220, 164, 3, 212, 111, 14, 99, 153, 50, 150, 150, 24, 59, 57, 249, 209, 210, 129, 24, 67, 211, 104, 92, 154, 91, 156, 174, 222, 115, 207, 176, 40, 33, 190, 72, 39, 159, 180, 104, 213, 211, 213, 75, 240, 142, 241, 18, 50, 217, 105, 243, 62, 8, 95, 155, 200, 211, 9, 82, 230, 1, 198, 221, 53, 208, 15, 191, 110, 242, 106, 79, 53, 230, 234, 224, 176, 236, 8, 86, 17, 215, 233, 172, 141, 165, 96, 219, 100, 243, 81, 114, 148, 108, 173, 105, 226, 51, 110, 61, 84, 163, 167, 7, 37, 175, 240, 49, 245, 67, 28, 96, 106, 40, 186, 96, 244, 249, 55, 222, 201, 113, 31, 146, 62, 24, 98, 151, 73, 187, 191, 13, 26, 166, 44, 252, 46, 193, 173, 72, 135, 78, 66, 202, 125, 57, 231, 119, 37, 184, 68, 217, 165, 236, 179, 36, 131, 38, 87, 110, 198, 198, 151, 2, 195, 203, 131, 205, 241, 129, 21, 145, 249, 243, 148, 171, 34, 234, 4, 143, 77, 144, 255, 216, 142, 174, 176, 241, 157, 198, 218, 82, 7, 56, 221, 120, 119, 208, 247, 141, 208, 100, 172, 203, 177, 88, 33, 129, 52, 44, 211, 36, 28, 8, 17, 99, 115, 19, 35, 47, 207, 98, 186, 129, 137, 152, 108, 102, 190, 61, 161, 197, 95, 84, 169, 102, 36, 67, 90, 85, 167, 200, 56, 22, 100, 2, 79, 29, 122, 101, 51, 112, 91, 205, 158, 252, 96, 215, 16, 87, 236, 139, 20, 181, 94, 248, 226, 6, 207, 85, 75, 187, 31, 193, 165, 23, 223, 147, 233, 66, 167, 42, 49, 244, 248, 136, 246, 4, 107, 125, 145, 161, 197, 198, 100, 229, 208, 52, 145, 6, 166, 128, 168, 108, 93, 111, 13, 126, 41, 204, 188, 211, 133, 230, 253, 218, 169, 17, 205, 239, 213, 153, 149, 153, 247, 86, 110, 130, 206, 14, 90, 190, 16, 179, 186, 128, 178, 64, 114, 5, 192, 238, 65, 77, 37, 208, 132, 136, 71, 77, 44, 51, 75, 108, 40, 135, 212, 120, 103, 113, 68, 239, 83, 30, 121, 184, 85, 207, 17, 42, 51, 221, 233, 249, 173, 170, 227, 210, 67, 143, 195, 97, 199, 17, 1, 253, 160, 224, 42, 217, 238, 11, 169, 202, 181, 181, 155, 152, 217, 232, 80, 127, 8, 213, 18, 25, 134, 99, 210, 80, 62, 223, 78, 161, 2, 69, 247, 124, 145, 204, 167, 33, 79, 131, 12, 5, 86, 192, 169, 65, 43, 60, 20, 199, 206, 106, 50, 36, 55, 207, 212, 177, 111, 37, 97, 168, 234, 37, 130, 50, 78, 120, 241, 125, 35, 222, 34, 203, 141, 126, 52, 92, 217, 37, 116, 29, 114, 234, 79, 85, 13, 114, 203, 66, 39, 252, 200, 192, 134, 100, 81, 149, 35, 219, 204, 112, 105, 210, 7, 118, 97, 66, 14, 94, 120, 24, 146, 61, 181, 222, 212, 240, 132, 114, 197, 177, 116, 74, 81, 39, 243, 74, 177, 207, 231, 76, 27, 247, 209, 67, 137, 79, 53, 145, 6, 166, 67, 199, 77, 5, 183, 219, 61, 151, 129, 23, 22, 187, 89, 93, 65, 89, 33, 153, 146, 128, 81, 166, 75, 161, 119, 117, 67, 101, 195, 33, 85, 253, 168, 196, 39, 4, 11, 211, 167, 168, 52, 56, 41, 106, 166, 98, 192, 34, 12, 241, 99, 229, 6, 229, 227, 69, 3, 198, 73, 107, 113, 30, 29, 239, 102, 74, 183, 78, 209, 68, 19, 166, 116, 246, 84, 208, 156, 11, 246, 152, 136, 219, 20, 169, 127, 60, 40, 119, 63, 123, 132, 39, 83, 151, 250, 212, 28, 76, 105, 58, 247, 138, 107, 207, 57, 55, 129, 45, 141, 216, 177, 77, 213, 192, 122, 111, 83, 83, 29, 63, 115, 67, 200, 179, 87, 26, 5, 152, 124, 219, 226, 250, 163, 158, 83, 238, 144, 3, 106, 186, 147, 65, 187, 219, 205, 162, 9, 103, 52, 212, 176, 108, 182, 171, 156, 26, 209, 155, 180, 219, 200, 231, 69, 31, 62, 254, 97, 34, 227, 21, 88, 7, 197, 43, 62, 66, 24, 225, 194, 75, 175, 220, 165, 65, 136, 198, 92, 220, 34, 143, 230, 192, 71, 103, 52, 96, 127, 37, 25, 53, 140, 173, 83, 150, 90, 102, 150, 93, 42, 141, 47, 249, 206, 109, 118, 12, 114, 104, 18, 120, 205, 29, 122, 131, 172, 107, 223, 214, 174, 232, 22, 247, 233, 155, 157, 54, 238, 216, 94, 96, 205, 253, 119, 96, 55, 212, 179, 186, 131, 178, 195, 218, 154, 89, 220, 235, 209, 222, 177, 68, 83, 198, 244, 77, 105, 232, 116, 193, 37, 157, 27, 218, 150, 78, 186, 234, 107, 38, 186, 0, 148, 17, 221, 120, 15, 65, 221, 75, 147, 30, 146, 186, 204, 56, 132, 77, 9, 211, 118, 169, 170, 124, 241, 111, 179, 193, 145, 105, 11, 33, 98, 82, 47, 21, 161, 39, 15, 247, 65, 127, 11, 91, 20, 156, 8, 168, 248, 100, 67, 59, 50, 190, 34, 149, 160, 83, 197, 154, 173, 35, 75, 176, 168, 118, 248, 229, 163, 210, 241, 1, 120, 81, 227, 246, 178, 139, 225, 20, 141, 245, 177, 19, 246, 131, 53, 83, 199, 126, 189, 114, 101, 103, 19, 157, 80, 184, 241, 149, 174, 26, 249, 158, 11, 49, 241, 144, 194, 155, 66, 35, 94, 243, 13, 225, 195, 147, 126, 145, 160, 100, 74, 80, 67, 85, 49, 189, 196, 128, 200, 164, 168, 241, 197, 80, 113, 170, 95, 211, 40, 246, 86, 232, 100, 104, 89, 43, 26, 152, 218, 101, 185, 15, 25, 245, 160, 83, 190, 104, 192, 112, 123, 162, 15, 11, 147, 78, 137, 148, 182, 74, 128, 147, 163, 119, 224, 220, 34, 150, 67, 119, 245, 254, 246, 142, 47, 249, 149, 244, 83, 144, 199, 16, 86, 180, 96, 205, 116, 245, 231, 57, 177, 150, 189, 252, 154, 244, 209, 73, 51, 161, 127, 243, 63, 248, 255, 219, 120, 83, 132, 67, 120, 123, 60, 250, 32, 177, 127, 56, 254, 184, 36, 238, 195, 63, 199, 79, 218, 143, 19, 252, 139, 14, 158, 181, 119, 183, 32, 221, 40, 86, 243, 88, 234, 41, 175, 12, 195, 151, 226, 67, 129, 49, 196, 94, 60, 86, 176, 214, 6, 119, 152, 73, 64, 179, 122, 36, 229, 72, 177, 217, 141, 34, 153, 51, 215, 15, 76, 0, 153, 8, 232, 100, 77, 90, 246, 156, 113, 65, 76, 168, 84, 176, 116, 215, 240, 132, 168, 71, 209, 52, 244, 176, 237, 25, 239, 244, 192, 83, 179, 225, 237, 162, 180, 51, 10, 229, 169, 188, 129, 217, 209, 212, 177, 15, 182, 45, 252, 135, 200, 71, 72, 18, 248, 5, 86, 143, 196, 200, 169, 138, 251, 83, 88, 156, 29, 109, 4, 240, 182, 145, 8, 18, 148, 248, 10, 246, 104, 209, 165, 18, 100, 82, 212, 163, 23, 108, 2, 24, 151, 168, 20, 92, 5, 78, 134, 86, 167, 184, 235, 158, 171, 62, 138, 81, 190, 238, 3, 204, 165, 239, 228, 232, 63, 238, 134, 11, 37, 54, 204, 84, 36, 209, 162, 235, 162, 229, 67, 243, 126, 121, 106, 88, 38, 236, 104, 225, 74, 174, 22, 90, 13, 168, 204, 2, 175, 10, 131, 211, 198, 125, 90, 250, 23, 252, 253, 221, 26, 89, 108, 8, 222, 30, 244, 75, 172, 40, 40, 230, 225, 139, 74, 111, 105, 5, 20, 80, 207, 109, 167, 29, 251, 26, 56, 125, 71, 44, 63, 31, 49, 163, 158, 160, 156, 80, 183, 29, 157, 158, 145, 81, 207, 22, 82, 170, 233, 182, 174, 34, 74, 179, 182, 255, 137, 185, 162, 91, 186, 167, 109, 97, 202, 52, 204, 205, 21, 230, 105, 183, 92, 75, 132, 180, 93, 219, 104, 96, 159, 40, 208, 52, 96, 108, 216, 45, 172, 189, 190, 99, 186, 47, 244, 123, 162, 191, 221, 190, 67, 210, 180, 212, 198, 116, 183, 148, 222, 120, 67, 199, 195, 73, 66, 79, 31, 210, 121, 123, 232, 170, 11, 36, 197, 127, 239, 245, 21, 181, 31, 219, 99, 62, 243, 223, 49, 253, 173, 80, 216, 178, 185, 15, 52, 239, 245, 53, 23, 183, 173, 70, 72, 143, 202, 203, 160, 189, 231, 75, 139, 131, 228, 187, 81, 143, 126, 234, 188, 210, 158, 6, 70, 222, 82, 110, 5, 201, 180, 253, 45, 11, 163, 242, 126, 191, 63, 190, 93, 196, 214, 157, 61, 244, 233, 26, 12, 103, 189, 100, 74, 18, 159, 58, 112, 238, 68, 93, 110, 122, 196, 12, 85, 211, 172, 226, 96, 16, 7, 243, 82, 211, 119, 234, 53, 59, 220, 254, 48, 98, 178, 187, 247, 56, 218, 148, 24, 9, 254, 69, 77, 229, 56, 84, 119, 234, 252, 98, 180, 147, 162, 56, 221, 152, 165, 30, 45, 21, 190, 199, 26, 184, 99, 199, 82, 196, 53, 221, 18, 85, 6, 190, 22, 68, 75, 43, 116, 211, 193, 179, 23, 89, 98, 87, 169, 60, 119, 59, 136, 22, 124, 10, 140, 67, 62, 179, 50, 59, 142, 206, 208, 194, 157, 26, 59, 88, 133, 92, 6, 222, 118, 33, 79, 161, 94, 160, 92, 144, 127, 10, 218, 111, 92, 117, 123, 44, 162, 253, 127, 62, 40, 233, 216, 107, 123, 253, 40, 45, 10, 130, 112, 164, 59, 88, 53, 85, 190, 235, 219, 226, 143, 87, 146, 81, 80, 12, 105, 120, 11, 152, 48, 38, 4, 44, 126, 219, 85, 110, 95, 127, 13, 160, 5, 41, 161, 141, 238, 74, 12, 72, 198, 217, 146, 160, 156, 141, 107, 188, 25, 246, 138, 122, 213, 166, 109, 254, 157, 154, 119, 53, 26, 60, 77, 156, 103, 1, 74, 162, 38, 130, 33, 109, 0, 150, 191, 8, 125, 170, 119, 86, 238, 188, 174, 71, 15, 235, 218, 39, 87, 42, 10, 245, 132, 26, 64, 184, 208, 128, 145, 92, 141, 69, 234, 125, 39, 94, 250, 183, 234, 139, 202, 139, 113, 141, 50, 53, 95, 89, 80, 4, 80, 114, 240, 108, 165, 220, 137, 161, 138, 10, 161, 210, 131, 250, 129, 242, 193, 6, 125, 202, 7, 185, 185, 234, 44, 23, 48, 114, 128, 83, 105, 160, 203, 133, 24, 120, 187, 254, 154, 8, 27, 82, 206, 0, 153, 54, 238, 251, 28, 23, 233, 165, 229, 89, 13, 66, 141, 101, 0}
)

func getCountryMapped() map[string]Country {
	// load + index countries into map
	// for below functions.
	onceCountry.Do(func() {
		rd := brotli.NewReader(bytes.NewBuffer(data))
		rb, _ := ioutil.ReadAll(rd)

		if err := json.Unmarshal(rb, &mapped); err != nil {
			panic(err.Error())
		}

		data = nil
	})

	return mapped
}

// GetCountry returns a single Country that matches the country
// code passed and whether it was found
func GetCountry(code string) (c Country, found bool) {
	c, found = getCountryMapped()[strings.ToLower(code)]
	return
}