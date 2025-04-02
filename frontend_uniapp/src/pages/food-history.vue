<template>
  <!-- No changes to template section -->
</template>

<script>
import { ref, onMounted } from 'vue';
import api from '@/api/food';

export default {
  setup() {
    const formData = ref({});
    const showEditModal = ref(false);

    // 编辑记录
    const handleEdit = (record) => {
      // 使用获取到的记录作为基础，只更新需要修改的字段
      const existingRecord = { ...record };
      // 删除不需要的字段
      delete existingRecord.created_at;
      delete existingRecord.updated_at;
      delete existingRecord.deleted_at;
      
      // 设置表单默认值
      formData.value = {
        ...existingRecord,
        record_time: new Date(record.record_time).toISOString().slice(0, 16) // 转换为本地时间格式
      };
      showEditModal.value = true;
    };

    // 保存编辑
    const saveEdit = async () => {
      try {
        if (!formData.value.food_name) {
          uni.showToast({
            title: '请输入食物名称',
            icon: 'none'
          });
          return;
        }

        // 确保 record_time 是 ISO 格式
        const updateData = {
          ...formData.value,
          record_time: new Date(formData.value.record_time).toISOString()
        };

        await api.updateFoodRecord(formData.value.id, updateData);
        uni.showToast({
          title: '更新成功',
          icon: 'success'
        });
        showEditModal.value = false;
        loadRecords();
      } catch (error) {
        console.error('更新记录失败:', error);
        uni.showToast({
          title: error.message || '更新失败',
          icon: 'none'
        });
      }
    };

    return {
      formData,
      showEditModal,
      handleEdit,
      saveEdit
    };
  }
};
</script>

<style>
  /* No changes to style section */
</style> 