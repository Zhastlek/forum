package pkg

func RemoveDuplicates(nums []int) []int {
	length := len(nums)
	if length <= 1 {
		return nums
	}
	var slice []int
	var status bool
	for i := 0; i < length; i++ {
		status = false
		for j := 0; j < len(slice); j++ {
			if nums[i] == slice[j] {
				status = true
			}
		}
		if !status {
			slice = append(slice, nums[i])
		}
	}
	return slice
}
